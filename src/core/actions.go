/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
*/

package Core

import (
    "os"
	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

var actions = []Types.Action {
    Shell,
    Copy,
    Move,
    Mkdir,
    Touch,
    Get,
}

func actionsToMap () map [ string ] Types.Action {
    cib := map [ string ] Types.Action {};
    for _, a := range actions {
        cib[ a.Name ] = a;
    }
    return cib;
}

var actList = actionsToMap();

func GetAction ( name string ) Types.Action {
    act, ok := actList[ name ];
    if !ok {
        Internal.ErrPrintf( "Fatal Error: Action with the name '%s' does not exist.\n", name );
        os.Exit( 1 );
    }
    return act;
}