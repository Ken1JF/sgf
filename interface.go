/*
 *  File:		src/github.com/Ken1JF/sgf/interface.go
 *  Project:	abst-hier
 *
 *  Created by Ken Friedenbach on 12/16/09.
 *  Copyright 2009-2014 Ken Friedenbach. All rights reserved.
 *
 *	This package implements the interface to the SGF Parser.
 *
 *	Much of this code is based on the file: src/pkg/go/parser/interface.go
 *	which is found in the Go Language project, and is used under the terms
 *	of the license cited below.
 *
 * Copyright 2009 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package sgf

import (
	"bytes"
	"github.com/Ken1JF/ah"
	"io"
	"io/ioutil"
)

// If src != nil, readSource converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
func readSource(filename string, src interface{}) ([]byte, ah.ErrorList) {
	var errs ah.ErrorList
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), errs
		case []byte:
			return s, errs
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.Bytes(), errs
			}
		case io.Reader:
			var buf bytes.Buffer
			_, err := io.Copy(&buf, s)
			if err != nil {
				errs.Add(ah.NoPos, err.Error())
				return nil, errs
			}
			return buf.Bytes(), errs
		default:
			errs.Add(ah.NoPos, "invalid source")
			return nil, errs
		}
	}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		errs.Add(ah.NoPos, err.Error())
	}
	return bytes, errs
}

// ParseFile parses a SGF file and returns a GoTree.TreeNode node.
//
// If src != nil, ParseFile parses the file source from src. src may
// be provided in a variety of formats. At the moment the following types
// are supported: string, []byte, and io.Reader. In this case, filename is
// only used for source position information and error messages.
//
// If src == nil, ParseFile parses the file specified by filename.
//
// The mode parameter controls the amount of source text parsed and other
// optional Parser functionality.
//
// If the source couldn't be read, the returned TreeNode is nil and the error
// indicates the specific failure. If the source was read but errors were found,
// the result is a partial TreeNode (with GoTree.BadX Nodes representing the fragments of erroneous source code). Multiple errors
// are returned via a Scanner.ErrorList which is sorted by file position.
//
// If the source couldn't be read, the returned TreeNode is nil and the error
// indicates the specific failure. If the source was read but syntax
// errors were found, the result is a partial tree (with TreeNode.BadX Nodes
// representing the fragments of erroneous SGF file). Multiple errors
// are returned via a Scanner.ErrorList which is sorted by file position.
func ParseFile(filename string, src interface{}, mode uint, moveLimit int) (*Parser, ah.ErrorList) {
	var p Parser
	var errL ah.ErrorList

	data, errL := readSource(filename, src)
	if len(errL) != 0 {
		return nil, errL
	}

	p.initParser(filename, data, mode, moveLimit)
	p.parseFile()
	return &p, p.errors
}
