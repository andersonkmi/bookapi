// Package repository contains the persistence layer of the Book API.
//
// It defines the GORM models stored in the database ([Book] and [Comment])
// together with [BookRepository], which groups the data-access operations used
// by the HTTP layer.
//
// Books are serialised to JSON with their comments flattened into an array of
// strings, so the wire format stays independent from the relational model:
//
//	{"id": 1, "title": "Dune", "author": "Frank Herbert", "comments": ["great"]}
//
// See [Book.MarshalJSON] and [Book.UnmarshalJSON] for the conversion rules.
//
// Errors are returned unwrapped from GORM, so callers can test for
// gorm.ErrRecordNotFound with errors.Is to detect missing records.
package repository
