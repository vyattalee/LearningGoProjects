/*
//构建约束以一行+build开始的注释。在+build之后列出了一些条件，在这些条件成立时，该文件应包含在编译的包中；
//约束可以出现在任何源文件中，不限于go文件；
//+build必须出现在package语句之前，+build注释之后应要有一个空行。
//多个条件之间，空格表示OR；逗号表示AND；叹号(!)表示NOT
//一个文件可以有多个+build，它们之间的关系是AND。
*/

// +build linux darwin

package os

const PathSeparator = "/"
