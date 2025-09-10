/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Cli

import (
	Types "github.com/neotesk/truct/src/types"
);

var Commands = []Types.Command{
    Init,
    Do,
    Setup,
    Run,
};

var CommandMap = UpdateCommandMap();

func UpdateCommandMap () map[ string ] Types.Command {
    output := map[ string ] Types.Command {};
    for _, command := range Commands {
        output[ command.Name ] = command;
    }
    return output;
}

func GetCommandByName ( name string ) *Types.Command {
    if command, ok := CommandMap[ name ]; ok {
        return &command;
    }
    return nil;
}