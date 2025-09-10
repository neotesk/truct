/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Types

type ProjectDetails struct {
    Name string
    Description string
    Version string
    Repository string
}

type TructSettings struct {
    Silent bool
    ReportActions bool
    Shell string
}

type WorkflowVariable struct {
    Name string
    Value string
    Optional bool
}

type DependencyList struct {
    FileDependencies []string
    CommandDependencies []string
}

type Workflow struct {
    Name string
    Description string
    Expects []WorkflowVariable
    Dependencies DependencyList
    Actions []any
}

type TructFileRaw struct {
    Project ProjectDetails
    Dependencies DependencyList
    ProjectRaw map[ string ] any
    Settings map[ string ] any
    Environment map[ string ] any
    Variables map[ string ] any
    Workflows map[ string ] Workflow
}

type TructFile struct {
    Project ProjectDetails
    Dependencies DependencyList
    ProjectRaw map[ string ] any
    Settings TructSettings
    Environment map[ string ] any
    Variables map[ string ] any
    Workflows map[ string ] Workflow
}

type CommandLineArgs struct {
    Arguments map[ string ] string;
    Flags map[ string ] bool;
    Keywords []string;
    TructFile TructFile;
}

type TructWorkflowRunArgs struct {
    WorkflowName string;
    ScopeVariables map[ string ] string;
    CommandLineArgs CommandLineArgs;
    WorkingDirectory string;
}