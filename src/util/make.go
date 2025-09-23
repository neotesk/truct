/*
   Truct, Pretty minimal workflow runner.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Internal

import (
    "os"
    "fmt"
    "maps"
    "reflect"
)

func SetField ( obj any, name string, value any ) error {
    structValue := reflect.ValueOf( obj ).Elem();
    fieldVal := structValue.FieldByName( name );
    if !fieldVal.IsValid() {
        return fmt.Errorf( "No such field: %s in obj", name );
    }
    if !fieldVal.CanSet() {
        return fmt.Errorf( "Cannot set %s field value", name );
    }
    val := reflect.ValueOf( value );
    if !val.IsValid() {
        return nil;
    }
    if fieldVal.Type() != val.Type() {
        if m,ok := value.( map[ string ] any ); ok {
            // if field value is struct
            if fieldVal.Kind() == reflect.Struct {
                return FillStruct( m, fieldVal.Addr().Interface() );
            }
            // if field value is a pointer to struct
            if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
                if fieldVal.IsNil() {
                    fieldVal.Set( reflect.New( fieldVal.Type().Elem() ) );
                }
                return FillStruct( m, fieldVal.Interface() );
            }
        }
        return fmt.Errorf( "Provided value type didn't match obj field type" );
    }
    fieldVal.Set( val );
    return nil;
}

func FillStruct ( m map[ string ] any, s any ) error {
    struc := reflect.ValueOf( s ).Elem().Type();
    strucKeys := struc.NumField();
    for i := range strucKeys {
        field := struc.Field( i );
        tag := field.Tag.Get( "import" );
        if tag == "" {
            tag = field.Name;
        }
        err := SetField( s, field.Name, m[ tag ] );
        if err != nil {
            return err;
        }
    }
    return nil;
}

func Make [ T any ] ( thing any ) T {
    output, ok := thing.( T );
    if !ok {
        ErrPrintf( "Fatal Error: Cannot convert object into desired type.\n" );
        os.Exit( 1 );
    }
    return output;
}

func MakeCoalesce [ T any ] ( thing any, def T ) T {
    if thing == nil {
        return def;
    }
    output, ok := thing.( T );
    if !ok {
        ErrPrintf( "Fatal Error: Cannot convert object into desired type.\n" );
        os.Exit( 1 );
    }
    return output;
}

func MakeArray [ T any ] ( thing any, def []T ) []T {
    if thing == nil {
        return def;
    }
    output, ok := thing.( []any );
    if !ok {
        ErrPrintf( "Fatal Error: Cannot convert object into desired type.\n" );
        os.Exit( 1 );
    }
    retval := []T {};
    for _, ret := range( output ) {
        typed, ok := ret.( T );
        if !ok {
            ErrPrintf( "Fatal Error: Cannot convert object into desired type.\n" );
            os.Exit( 1 );
        }
        retval = append( retval, typed );
    }
    return retval;
}

func MakeSure ( thing any, errMsg string ) {
    if thing == nil {
        fmt.Fprintln( os.Stderr, "Fatal Error: " + errMsg );
        os.Exit( 1 );
    }
}

func PossibleItem [ T any ] ( arr []T, idx int ) any {
    if len( arr ) <= idx {
        return nil;
    }
    return arr[ idx ];
}

func CopyMap ( thing map [ string ] any ) map [ string ] any {
    targetMap := make( map [string] any );
    maps.Copy( targetMap, thing );
    return targetMap;
}