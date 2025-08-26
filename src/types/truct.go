/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
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

type Workflow struct {
    Name string
    Description string
    Expects []WorkflowVariable
    Actions []any
}

type TructFileRaw struct {
    Project ProjectDetails
    Settings map[ string ] any
    Environment map[ string ] any
    Variables map[ string ] any
    Workflows map[ string ] Workflow
}

type TructFile struct {
    Project ProjectDetails
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
}