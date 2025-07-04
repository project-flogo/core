package debugger

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/engine/support"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	"os"
	"path"
	"reflect"
	"strings"
)

func GenerateReport(config *trigger.HandlerConfig, interceptors []*support.TaskInterceptor, coverage *support.Coverage, instanceID string, flowInputs map[string]interface{}, flowOutputs map[string]interface{}) {
	finalReport := &support.OutputReport{
		AppName:    GetAppName(),
		AppVersion: GetAppVersion(),
		InstanceID: instanceID,
	}
	report := &support.Report{}

	triggerNode := &support.Trigger{
		ID:       config.Parent.Id,
		Settings: config.Settings,
	}

	handler := support.Handler{
		FlowName: config.Name,
		Input:    flowInputs,
		Output:   flowOutputs,
	}

	triggerNode.Handler = handler

	report.Trigger = triggerNode

	report.Flows = processFlowReport(config.Name, interceptors, coverage)

	fileName := config.Name + "-" + instanceID + ".json"

	finalReport.Report = report
	op, err := json.MarshalIndent(finalReport, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling report ", err)
	}

	reportPath := os.Getenv("FLOW_EXECUTION_FILES")

	if reportPath == "" {
		reportPath = path.Join(os.TempDir(), "flow-executions", GetAppName(), fileName)
	} else {
		reportPath = path.Join(reportPath, "flow-executions", GetAppName(), fileName)
	}

	log.RootLogger().Infof("Generate Report for Flow Execution: %s", reportPath)

	os.MkdirAll(path.Dir(reportPath), os.ModePerm)
	err = os.WriteFile(reportPath, op, 0777)
	if err != nil {
		fmt.Printf("Error writing report to file: %v", err)
	}
}

func processFlowReport(mainFlow string, interceptors []*support.TaskInterceptor, coverage *support.Coverage) *support.FlowReport {

	dataMap := getSubFlowDataMap(mainFlow, coverage)

	subFlowMap := make(map[string]map[string]string)
	for _, subFlowCoverage := range coverage.SubFlowCoverage {
		if val, ok := subFlowMap[subFlowCoverage.HostFlow]; ok {
			val[subFlowCoverage.SubFlowActivity] = subFlowCoverage.SubFlowName
		} else {
			subFlowMap[subFlowCoverage.HostFlow] = make(map[string]string)
			subFlowMap[subFlowCoverage.HostFlow][subFlowCoverage.SubFlowActivity] = subFlowCoverage.SubFlowName
		}

	}

	flowReport := &support.FlowReport{}
	flowReport.Name = mainFlow
	var flowOpReport *support.ActivityReport
	var activityReportList []support.ActivityReport
	var errorHandlerActivityReport = make([]support.ActivityReport, 0)
	var linkReportList []support.LinkReport
	var errorHandlerLinkReportList = make([]support.LinkReport, 0)

	var interceptorMap = make(map[string]*support.TaskInterceptor)
	for _, interceptor := range interceptors {
		activityName := strings.Replace(interceptor.ID, mainFlow+"-", "", 1)
		interceptorMap[activityName] = interceptor
	}

	for _, activity := range coverage.ActivityCoverage {

		if activity.FlowName != mainFlow {
			continue
		}
		activityReport := &support.ActivityReport{}
		activityReport.ActivityName = activity.ActivityName
		activityReport.Inputs = activity.Inputs
		if !reflect.ValueOf(&activity.Outputs).IsNil() {
			activityReport.Outputs = &activity.Outputs
		} else {
			activityReport.Outputs = nil
		}

		activityReport.Error = activity.Error

		if activity.IsMainFlow {
			activityReportList = append(activityReportList, *activityReport)
		} else {
			errorHandlerActivityReport = append(errorHandlerActivityReport, *activityReport)
		}

	}

	for _, link := range coverage.TransitionCoverage {

		if link.FlowName != mainFlow {
			continue
		}
		linkReport := &support.LinkReport{}
		linkReport.LinkName = link.TransitionName
		linkReport.To = link.TransitionTo
		linkReport.From = link.TransitionFrom
		if link.IsMainFlow {
			linkReportList = append(linkReportList, *linkReport)
		} else {
			errorHandlerLinkReportList = append(errorHandlerLinkReportList, *linkReport)
		}

	}

	_, ok := interceptorMap["_flowOutput"]
	if ok {
		flowOpReport = &support.ActivityReport{}
		flowOpReport.ActivityName = "FlowReport Output"

	}

	if flowOpReport != nil {
		activityReportList = append(activityReportList, *flowOpReport)
	}
	flowReport.ActivityReport = activityReportList
	flowReport.FlowErrorHandler = support.FlowErrorHandler{
		ActivityReport: errorHandlerActivityReport,
		LinkReport:     errorHandlerLinkReportList,
	}
	flowReport.LinkReport = linkReportList

	for k, testReport := range dataMap {
		if val, ok := subFlowMap[k]; ok {
			subMap := make(map[string]interface{})
			for k1, v1 := range val {
				if val, ok := dataMap[v1]; ok {
					subMap[k1] = val
				}
			}
			testReport.SubFlow = subMap
		}
	}

	if _, ok := subFlowMap[mainFlow]; ok {
		subFlow := subFlowMap[mainFlow]
		subMap := make(map[string]interface{})
		for k, v := range subFlow {
			if val, ok := dataMap[v]; ok {
				subMap[k] = val
			}
		}
		flowReport.SubFlow = subMap
	}

	return flowReport
}

func getSubFlowDataMap(mainFlow string, coverage *support.Coverage) map[string]*support.FlowReport {
	subFlowList := make(map[string]*support.FlowReport)

	for _, activity := range coverage.ActivityCoverage {
		if activity == nil {
			continue
		}
		if activity.FlowName == mainFlow {
			continue
		}
		activityReport := &support.ActivityReport{}
		activityReport.ActivityName = activity.ActivityName
		activityReport.Inputs = activity.Inputs
		activityReport.Outputs = &activity.Outputs
		activityReport.Error = activity.Error

		val, ok := subFlowList[activity.FlowName]
		if ok {
			if activity.IsMainFlow {
				val.ActivityReport = append(val.ActivityReport, *activityReport)
			} else {
				val.FlowErrorHandler.ActivityReport = append(val.FlowErrorHandler.ActivityReport, *activityReport)
			}
		} else {
			val = &support.FlowReport{
				Name:           activity.FlowName,
				ActivityReport: make([]support.ActivityReport, 0),
				LinkReport:     make([]support.LinkReport, 0),
				FlowErrorHandler: support.FlowErrorHandler{
					ActivityReport: make([]support.ActivityReport, 0),
					LinkReport:     make([]support.LinkReport, 0),
				},
				SubFlow: make(map[string]interface{}),
			}
			if activity.IsMainFlow {
				val.ActivityReport = append(val.ActivityReport, *activityReport)
			} else {
				val.FlowErrorHandler.ActivityReport = append(val.FlowErrorHandler.ActivityReport, *activityReport)
			}
			subFlowList[activity.FlowName] = val
		}
	}

	for _, link := range coverage.TransitionCoverage {

		if link.FlowName == mainFlow {
			continue
		}
		linkReport := &support.LinkReport{}
		linkReport.LinkName = link.TransitionName
		linkReport.To = link.TransitionTo
		linkReport.From = link.TransitionFrom

		val, ok := subFlowList[link.FlowName]
		if ok {
			if link.IsMainFlow {
				val.LinkReport = append(val.LinkReport, *linkReport)
			} else {
				val.FlowErrorHandler.LinkReport = append(val.FlowErrorHandler.LinkReport, *linkReport)
			}
		}
	}

	return subFlowList

}
