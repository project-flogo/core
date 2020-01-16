 
<p align="center">
  <img src ="https://raw.githubusercontent.com/TIBCOSoftware/flogo/master/images/projectflogo.png" />
</p>

<p align="center" >
  <b>Serverless functions and edge microservices made painless</b>
</p>

<p align="center">
  <img src="https://travis-ci.org/project-flogo/core.svg?branch=master"/>
  <img src="https://img.shields.io/badge/dependencies-up%20to%20date-green.svg"/>
  <img src="https://img.shields.io/badge/license-BSD%20style-blue.svg"/>
  <a href="https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link"><img src="https://badges.gitter.im/Join%20Chat.svg"/></a>
</p>

<p align="center">
  <a href="#getting-started">Getting Started</a> | <a href="#documentation">Documentation</a> | <a href="#contributing">Contributing</a> | <a href="#license">License</a>
</p>

<br/>

Project Flogo is an open source framework to simplify building efficient & modern serverless functions and edge microservices and _this_ repository is the core library used to create and extend those **Flogo Applications**. 

# Flogo Core
Flogo Core is the core flogo library which contains the apis to create and extend Flogo applications.

## Getting started
If you want to get started with [Project Flogo](flogo.io), you should install the install the [Flogo CLI](https://github.com/project-flogo/cli).  You can find details there on creating a quick sample application.  You also might want to check out the [getting started](https://tibcosoftware.github.io/flogo/getting-started/) guide in our docs or check out the [Labs](https://tibcosoftware.github.io/flogo/labs/) section in our docs for in depth tutorials.

## Documentation
Here is some documentation to help you get started understanding some of the fundamentals of the Flogo Core library. 

* [Model](docs/model.md): The Flogo application model
* [Data Types](docs/datatypes.md): The Flogo data types
* [Mapping](docs/mapping.md): Mapping data in Flogo applications

In addition to low-level APIs used to support and run Flogo applications, the Core library contains some high-level APIs.  There is an API that can be used to programmatically create and run an application.  There are also interfaces that can be implemented to create your own Flogo contributions, such as Triggers and Activities. 

* [Application](docs/app-api.md): API to build and execute a Flogo application
* [Contributions](docs/contribs.md): APIs and interfaces for Flogo contribution development

## Contributing
Want to contribute to Project Flogo? We've made it easy, all you need to do is fork the repository you intend to contribute to, make your changes and create a Pull Request! Once the pull request has been created, you'll be prompted to sign the CLA (Contributor License Agreement) online.

Not sure where to start? No problem, you can browse the Project Flogo repos and look for issues tagged `kind/help-wanted` or `good first issue`. To make this even easier, we've added the links right here too!
* Project Flogo: [kind/help-wanted](https://github.com/TIBCOSoftware/flogo/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/TIBCOSoftware/flogo/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
* flogo cli: [kind/help-wanted](https://github.com/project-flogo/cli/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/project-flogo/cli/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
* flogo core: [kind/help-wanted](https://github.com/project-flogo/core/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/project-flogo/core/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
* flogo contrib: [kind/help-wanted](https://github.com/project-flogo/contrib/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) and [good first issue](https://github.com/project-flogo/contrib/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

Another great way to contribute to Project Flogo is to check [flogo-contrib](https://github.com/project-flogo/contrib). That repository contains some basic contributions, such as activities, triggers, etc. Perhaps there is something missing? Create a new activity or trigger or fix a bug in an existing activity or trigger.

If you have any questions, feel free to post an issue and tag it as a question, email flogo-oss@tibco.com or chat with the team and community:

* The [project-flogo/Lobby](https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for general discussions, start here for all things Flogo!
* The [project-flogo/developers](https://gitter.im/project-flogo/developers?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for developer/contributor focused conversations. 

For additional details, refer to the [Contribution Guidelines](https://github.com/TIBCOSoftware/flogo/blob/master/CONTRIBUTING.md).

## License 
Flogo source code in [this](https://github.com/project-flogo/core) repository is under a BSD-style license, refer to [LICENSE](https://github.com/project-flogo/core/blob/master/LICENSE) 
