/*
    Truct, Pretty minimal workflow runner.
    Open-Source, WTFPL License.

    Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package main

import (
    "os"
    "strings"
    Types "github.com/neotesk/truct/src/types"
    Internal "github.com/neotesk/truct/src/util"
)

func Arguments ( defaults Types.DefaultArgs ) Types.ArgumentsList {
    // Get arguments list
    args := os.Args[ 1: ];

    // Create a list for arguments, flags and keywords
    // so we can use them later on.
    arguments := make( map[ string ] string );
    flags := map[ string ] bool {};
    keywords := [] string{};

    // Iterate over defaults
    for _, arg := range defaults.Arguments {
        arguments[ arg.Name ] = arg.DefaultValue
    }
    for _, arg := range defaults.Flags {
        flags[ arg.Name ] = arg.DefaultValue
    }

    // Create an iterator
    currentIdx := 0;
    currentItem := "";
    argsLen := len( args );
    for currentIdx < argsLen {
        currentItem = args[ currentIdx ];
        if strings.HasPrefix( currentItem, "--" ) {
            // This is an argument, We should check if
            // it contains an equals sign and if not we,
            // must expect a string at the next index.
            signIndex := strings.Index( currentItem, "=" );
            possibleName := currentItem[ 2: ];
            if signIndex != -1 {
                // This means that it exists! we should
                // now use it.
                items := strings.Split( possibleName, "=" );
                arguments[ items[ 0 ] ] = items[ 1 ];
            } else if currentIdx < argsLen - 1 {
                // This means that we have one more
                // argument inside our args list
                currentIdx++;
                arguments[ possibleName ] = args[ currentIdx ];
            } else {
                // Whoops, seems like there's an issue and this
                // MUST NOT happen. We are missing the value so
                // literally print an error and exit
                Internal.ErrPrintf( "Fatal Error! Given argument with the name '%s' is missing a value.\n", possibleName );
                os.Exit( 127 );
            }
        } else if strings.HasPrefix( currentItem, "-" ) {
            // This is a flag, since we allow flag
            // combinations, append every char.
            for _, char := range currentItem[ 1: ] {
                flags[ string( char ) ] = true;
            }
        } else {
            // At this point this is a keyword
            keywords = append( keywords, currentItem );
        }
        currentIdx++;
    }

    return Types.ArgumentsList{
        Arguments: arguments,
        Flags: flags,
        Keywords: keywords,
    }
}