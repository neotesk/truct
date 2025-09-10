/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Internal

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type COLOR_LEVEL int;

const (
    C_UNKNOWN COLOR_LEVEL = iota;
    C_NONE;
    C_BASIC;
    C_COLOR256;
    C_TRUECOLOR;
);

func hueToRgb ( p, q, t float64 ) float64 {
    if t < 0 { t += 1; }
    if t > 1 { t -= 1; }
    if t < ( 1.0 / 6.0 ) { return p + ( q - p ) * 6.0 * t; }
    if t < ( 1.0 / 2.0 ) { return q; }
    if t < ( 2.0 / 3.0 ) { return p + ( q - p ) * ( 2.0 / 3.0 - t ) * 6.0; }
    return p;
}

func envMap () map [ string ] string {
    env := os.Environ();
    mp := map [ string ] string {}
    for _, envItem := range env {
        idx := strings.Index( envItem, "=" );
        mp[ envItem[ :idx ] ] = envItem[ idx + 1: ]
    }
    return mp;
}

type IColor struct {};
var Color IColor;

type RGB struct {
    R uint8
    G uint8
    B uint8
};

func ( IColor ) FromHEX ( hexString string ) ( *RGB, error ) {
    var rgb RGB;
	values, err := strconv.ParseUint( string( hexString ), 16, 32 );
	if err != nil {
		return &RGB {}, err;
	}
	rgb = RGB {
		R: uint8( values >> 16 ),
		G: uint8( ( values >> 8 ) & 0xFF ),
		B: uint8( values & 0xFF ),
	};
	return &rgb, nil;
}

func ( IColor ) FromRGB ( r, g, b uint8 ) RGB {
    return RGB { R: r, G: g, B: b };
}

func ( IColor ) FromHSL ( h, s, l float64 ) RGB {
    if s == 0 {
        gr := uint8( math.Round( 255.0 * l ) );
        return RGB { R: gr, G: gr, B: gr };
    }
    q := 0.0; if l < 0.5 { q = l * ( 1.0 + s ); } else { q = l + s - l * s; }
    p := 2 * l - q;
    x := 1.0 / 3.0;
    return RGB {
        R: uint8( math.Round( 255.0 * hueToRgb( p, q, h + x ) ) ),
        G: uint8( math.Round( 255.0 * hueToRgb( p, q, h ) ) ),
        B: uint8( math.Round( 255.0 * hueToRgb( p, q, h - x ) ) ),
    }
}

func wrapString ( col, text string ) string {
    return col + text + "\x1b[0m";
}

type IColorAdapter struct {};
var ColorAdapter IColorAdapter;

type KeyedColor struct {
    color RGB;
    key uint8;
}

var colorMap16 = []KeyedColor {
    { color: Color.FromRGB( 0  , 0  , 0   ), key: 30 },
    { color: Color.FromRGB( 128, 0  , 0   ), key: 31 },
    { color: Color.FromRGB( 0  , 128, 0   ), key: 32 },
    { color: Color.FromRGB( 128, 128, 0   ), key: 33 },
    { color: Color.FromRGB( 0  , 0  , 128 ), key: 34 },
    { color: Color.FromRGB( 255, 0  , 128 ), key: 35 },
    { color: Color.FromRGB( 0  , 128, 128 ), key: 36 },
    { color: Color.FromRGB( 128, 128, 128 ), key: 37 },
    { color: Color.FromRGB( 64 , 64 , 64  ), key: 90 },
    { color: Color.FromRGB( 255, 64 , 64  ), key: 91 },
    { color: Color.FromRGB( 64 , 255, 64  ), key: 92 },
    { color: Color.FromRGB( 255, 255, 64  ), key: 93 },
    { color: Color.FromRGB( 64 , 64 , 255 ), key: 94 },
    { color: Color.FromRGB( 255, 64 , 255 ), key: 95 },
    { color: Color.FromRGB( 64 , 255, 255 ), key: 96 },
    { color: Color.FromRGB( 255, 255, 255 ), key: 97 },
};

func ( IColorAdapter ) To256Color ( col RGB, isBG bool ) string {
    cSwitch := 3; if isBG { cSwitch = 4; }
    format := "\x1b[" + strconv.Itoa( cSwitch ) + "8;5;%sm";
    if col.R == col.G && col.G == col.B {
        if col.R < 8 { return fmt.Sprintf( format, strconv.Itoa( 16 ) ); }
        if col.R > 248 { return fmt.Sprintf( format, strconv.Itoa( 231 ) ); }
        return fmt.Sprintf( format, strconv.Itoa( int( math.Round( ( ( float64( col.R ) - 8.0 ) / 247.0 ) * 24.0 ) + 232 ) ) );
    }
    return fmt.Sprintf( format, strconv.Itoa( 16 +
        int( 36 * math.Round( float64( col.R ) / 255 * 5 ) ) +
        int( 6  * math.Round( float64( col.G ) / 255 * 5 ) ) +
        int( math.Round( float64( col.B ) / 255 * 5 ) ) ) );
}

func ( IColorAdapter ) ToTrueColor ( col RGB, isBG bool ) string {
    cSwitch := 3; if isBG { cSwitch = 4; }
    return fmt.Sprintf( "\x1b[%s8;2;%s;%s;%sm",
        strconv.Itoa( cSwitch ),
        strconv.Itoa( int( col.R ) ),
        strconv.Itoa( int( col.G ) ),
        strconv.Itoa( int( col.B ) ) );
}

func ( IColorAdapter ) To16Color ( col RGB, isBG bool ) string {
    // We are going to calculate the closest color
    // distance in our map, then we will get the color
    // ID from the map. For calculating the closest color,
    // we will be using Euclidean Distance formula
    clDistance := math.Inf( 1 );
    clColor := uint8( 0 );

    R := col.R;
    G := col.G;
    B := col.B;

    for _, item := range colorMap16 {
        distance := math.Sqrt(
            math.Pow( float64( R - item.color.R ), 2 ) +
            math.Pow( float64( G - item.color.G ), 2 ) +
            math.Pow( float64( B - item.color.B ), 2 ) );
        if distance < clDistance {
            clDistance = distance;
            clColor = item.key
        }
    }

    output := int( clColor );
    if isBG {
        output = output + 10;
    }

    return fmt.Sprintf( "\x1b[%sm", strconv.Itoa( output ) );
}

var env = envMap();

func In [ T any ] ( mp map[ string ] T, key... string ) bool {
    for _, v := range key {
        if _, ok := mp[ v ]; ok { return true; }
    }
    return false;
}

func envColor () COLOR_LEVEL {
    if !In( env, "COLOR" ) {
        return C_UNKNOWN;
    }
    item := env[ "COLOR" ];
    if item == "0" { return C_NONE; }
    if item == "1" { return C_BASIC; }
    if item == "2" { return C_COLOR256; }
    if item == "3" { return C_TRUECOLOR; }
    return C_UNKNOWN;
}

func supportedMaxColor () COLOR_LEVEL {
    forced := envColor();
    term := "dumb";
    if In( env, "TERM" ) {
        term = env[ "TERM" ];
    }
    if ( forced != C_UNKNOWN ) {
        return forced;
    }
    if ( In( env, "TF_BUILD" ) && In( env, "AGENT_NAME" ) ) || runtime.GOOS == "windows" {
        return C_BASIC;
    }
    if In( env, "ZED_ENVIRONMENT" ) || In( env, "ZED_TERM" ) || term == "xterm-kitty" || term == "xterm-gnome" || term == "alacritty" {
        return C_TRUECOLOR;
    }
    if term == "dumb" {
        return C_NONE;
    }
    if In( env, "CI" ) {
        if In( env, "GITHUB_ACTIONS" ) || In( env, "GITEA_ACTIONS" ) {
            return C_TRUECOLOR;
        }
        if In( env, "TRAVIS", "CIRCLECI", "APPVEYOR", "GITLAB_CI", "BUILDKITE","DRONE" ) {
            return C_BASIC;
        }
    }
    if strings.Contains( term, "-256" ) {
        return C_COLOR256;
    }
    if  strings.HasPrefix( term, "screen" ) ||
        strings.HasPrefix( term, "xterm" ) ||
        strings.HasPrefix( term, "vt100" ) ||
        strings.HasPrefix( term, "vt220" ) ||
        strings.HasPrefix( term, "rxvt" ) ||
        strings.Contains( term, "color" ) ||
        strings.Contains( term, "ansi" ) ||
        term == "cygwin" || term == "linux" {
        return C_BASIC;
    }
    return C_NONE;
}

var maxColor = supportedMaxColor();

func emptyColor ( r RGB, x bool ) string {
    return ""
}

func getColorer () func ( RGB, bool ) string {
    if maxColor == C_BASIC {
        return ColorAdapter.To16Color;
    }
    if maxColor == C_COLOR256 {
        return ColorAdapter.To256Color;
    }
    if maxColor == C_TRUECOLOR {
        return ColorAdapter.ToTrueColor;
    }
    return emptyColor;
}

var colorer = getColorer();

func color ( hex string ) string {
    rgb := HandleError( Color.FromHEX( hex ) );
    return colorer( *rgb, false );
}

func bgColor ( hex string ) string {
    rgb := HandleError( Color.FromHEX( hex ) );
    return colorer( *rgb, true );
}

func Colorify ( text, hex string ) string {
    return wrapString( color( hex ), text );
}

func ColorifyBG ( text, hex string ) string {
    return wrapString( bgColor( hex ), text );
}

func Boldify ( text string ) string {
    return "\x1b[1m" + text + "\x1b[22m";
}

var IsDebug = false;

func ErrPrintf ( format string, args... any ) {
    fmt.Fprintf( os.Stderr, Boldify( Colorify( format, "ff4040" ) ), args... );
    if IsDebug {
        panic( fmt.Sprintf( format, args... ) );
    }
}