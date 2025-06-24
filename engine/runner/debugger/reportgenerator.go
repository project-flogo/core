package debugger

import (
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/engine/runner/types"
	"github.com/project-flogo/core/trigger"
	"io/ioutil"
	"os"
	"strings"
)

func GenerateReport(config *trigger.HandlerConfig, interceptors []*types.TaskInterceptor, coverage *types.Coverage, instanceID string, flowInputs map[string]interface{}, flowOutputs map[string]interface{}) {
	finalReport := &types.OutputReport{}
	report := &types.Report{}

	ref := config.Parent.Ref

	trigger := &types.Trigger{
		ID:       config.Parent.Id,
		Settings: config.Settings,
	}

	handler := types.Handler{
		FlowName: config.Name,
		Input:    flowInputs,
		Output:   flowOutputs,
	}

	trigger.Handler = handler

	report.Trigger = trigger

	report.Flows = processFlowReport(config.Name, interceptors, coverage)

	fileName := ""
	if ref == "#startuphook" {
		fileName = "OnStartup-" + config.Name + ".json"
	} else if ref == "#shutdownhook" {
		fileName = "OnShutdown-" + config.Name + ".json"
	} else {
		fileName = config.Name + "-" + instanceID + ".json"
	}
	finalReport.Report = report
	op, err := json.MarshalIndent(finalReport, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling report ", err)
	}

	os.Remove(fileName)
	err = ioutil.WriteFile(fileName, op, 777)

}

func processFlowReport(mainFlow string, interceptors []*types.TaskInterceptor, coverage *types.Coverage) *types.FlowReport {

	dataMap := getSubFlowDataMap(mainFlow, interceptors, coverage)

	subFlowMap := make(map[string]map[string]string)
	for _, subFlowCoverage := range coverage.SubFlowCoverage {
		if val, ok := subFlowMap[subFlowCoverage.HostFlow]; ok {
			val[subFlowCoverage.SubFlowActivity] = subFlowCoverage.SubFlowName
		} else {
			subFlowMap[subFlowCoverage.HostFlow] = make(map[string]string)
			subFlowMap[subFlowCoverage.HostFlow][subFlowCoverage.SubFlowActivity] = subFlowCoverage.SubFlowName
		}

	}

	flowReport := &types.FlowReport{}
	flowReport.Name = mainFlow
	var flowOpReport *types.ActivityReport
	var activityReportList []types.ActivityReport
	var errorHandlerActivityReport []types.ActivityReport = make([]types.ActivityReport, 0)
	var linkReportList []types.LinkReport
	var errorHandlerLinkReportList []types.LinkReport = make([]types.LinkReport, 0)

	var interceptorMap = make(map[string]*types.TaskInterceptor)
	for _, interceptor := range interceptors {
		activityName := strings.Replace(interceptor.ID, mainFlow+"-", "", 1)
		interceptorMap[activityName] = interceptor
	}

	for _, activity := range coverage.ActivityCoverage {

		if activity.FlowName != mainFlow {
			continue
		}
		activityReport := &types.ActivityReport{}
		activityReport.ActivityName = activity.ActivityName
		activityReport.Inputs = activity.Inputs
		activityReport.Outputs = &activity.Outputs
		activityReport.Error = activity.Error

		if activity != nil {

			if activity.IsMainFlow {
				activityReportList = append(activityReportList, *activityReport)
			} else {
				errorHandlerActivityReport = append(errorHandlerActivityReport, *activityReport)
			}
		} else {
			activityReportList = append(activityReportList, *activityReport)
		}
	}

	for _, link := range coverage.TransitionCoverage {

		if link.FlowName != mainFlow {
			continue
		}
		linkReport := &types.LinkReport{}
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
		flowOpReport = &types.ActivityReport{}
		flowOpReport.ActivityName = "FlowReport Output"

	}

	if flowOpReport != nil {
		activityReportList = append(activityReportList, *flowOpReport)
	}
	flowReport.ActivityReport = activityReportList
	flowReport.FlowErrorHandler = types.FlowErrorHandler{
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

func getSubFlowDataMap(mainFlow string, interceptors []*types.TaskInterceptor, coverage *types.Coverage) map[string]*types.FlowReport {
	subFlowList := make(map[string]*types.FlowReport)

	for _, activity := range coverage.ActivityCoverage {
		if activity == nil {
			continue
		}
		if activity.FlowName == mainFlow {
			continue
		}
		activityReport := &types.ActivityReport{}
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
			val = &types.FlowReport{
				Name:           activity.FlowName,
				ActivityReport: make([]types.ActivityReport, 0),
				LinkReport:     make([]types.LinkReport, 0),
				FlowErrorHandler: types.FlowErrorHandler{
					ActivityReport: make([]types.ActivityReport, 0),
					LinkReport:     make([]types.LinkReport, 0),
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
		linkReport := &types.LinkReport{}
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
