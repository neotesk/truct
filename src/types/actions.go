/*
    Truct, Pretty minimal workflow manager.
    Open-Source, WTFPL License.

    Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Types

type ActionFunction func ( cwd string, details map[ string ] any, custParams TructWorkflowRunArgs ) error;

type Action struct {
    Name string;
    Expects map[ string ] any;
    Action ActionFunction;
}