package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

// handlerArgs provides typed access to MCP tool arguments.
type handlerArgs map[string]interface{}

func getArgs(request mcp.CallToolRequest) (handlerArgs, *mcp.CallToolResult) {
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, mcp.NewToolResultError("invalid arguments")
	}
	return handlerArgs(args), nil
}

func (a handlerArgs) getString(key string) (string, *mcp.CallToolResult) {
	v, ok := a[key].(string)
	if !ok {
		return "", mcp.NewToolResultError(fmt.Sprintf("%s must be a string", key))
	}
	return v, nil
}

func (a handlerArgs) getFloat(key string) (float64, *mcp.CallToolResult) {
	v, ok := a[key].(float64)
	if !ok {
		return 0, mcp.NewToolResultError(fmt.Sprintf("%s must be a number", key))
	}
	return v, nil
}

func (a handlerArgs) getOptionalFloat(key string, defaultVal float64) float64 {
	v, ok := a[key].(float64)
	if !ok || v <= 0 {
		return defaultVal
	}
	return v
}

func marshalResult(v interface{}) (*mcp.CallToolResult, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func handleGetBooks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	books, err := bookService.GetAllBooks()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get books: %v", err)), nil
	}
	return marshalResult(books)
}

func handleGetBook(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, errResult := getArgs(request)
	if errResult != nil {
		return errResult, nil
	}

	identifier, errResult := args.getString("identifier")
	if errResult != nil {
		return errResult, nil
	}

	book, err := bookService.GetBookOrBySlug(identifier)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get book: %v", err)), nil
	}
	if book == nil {
		return mcp.NewToolResultError("Book not found"), nil
	}

	return marshalResult(book)
}

func handleGetChapters(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, errResult := getArgs(request)
	if errResult != nil {
		return errResult, nil
	}

	bookIDStr, errResult := args.getString("book_id")
	if errResult != nil {
		return errResult, nil
	}

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		return mcp.NewToolResultError("Invalid book_id"), nil
	}

	chapters, err := chapterService.GetChaptersByBook(uint(bookID))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get chapters: %v", err)), nil
	}

	return marshalResult(chapters)
}

func handleGetChapterHadiths(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, errResult := getArgs(request)
	if errResult != nil {
		return errResult, nil
	}

	chapterID, errResult := args.getFloat("chapter_id")
	if errResult != nil {
		return errResult, nil
	}

	page := int(args.getOptionalFloat("page", 1))
	limit := int(args.getOptionalFloat("limit", 20))
	if limit > 100 {
		limit = 100
	}

	hadiths, total, err := hadithService.GetHadithsByChapter(uint(chapterID), page, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get hadiths: %v", err)), nil
	}

	return marshalResult(map[string]interface{}{
		"data": hadiths,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func handleGetHadith(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, errResult := getArgs(request)
	if errResult != nil {
		return errResult, nil
	}

	id, errResult := args.getFloat("id")
	if errResult != nil {
		return errResult, nil
	}

	hadith, err := hadithService.GetHadith(uint(id))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get hadith: %v", err)), nil
	}
	if hadith == nil {
		return mcp.NewToolResultError("Hadith not found"), nil
	}

	return marshalResult(hadith)
}

func handleGetHadithByNumber(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, errResult := getArgs(request)
	if errResult != nil {
		return errResult, nil
	}

	bookID, errResult := args.getFloat("book_id")
	if errResult != nil {
		return errResult, nil
	}

	number, errResult := args.getFloat("number")
	if errResult != nil {
		return errResult, nil
	}

	hadith, err := hadithService.GetHadithByBookAndNumber(uint(bookID), int(number))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get hadith: %v", err)), nil
	}
	if hadith == nil {
		return mcp.NewToolResultError("Hadith not found"), nil
	}

	return marshalResult(hadith)
}

func handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, errResult := getArgs(request)
	if errResult != nil {
		return errResult, nil
	}

	query, errResult := args.getString("query")
	if errResult != nil || query == "" {
		return mcp.NewToolResultError("query must be a non-empty string"), nil
	}

	page := int(args.getOptionalFloat("page", 1))
	limit := int(args.getOptionalFloat("limit", 20))
	if limit > 100 {
		limit = 100
	}

	results, total, err := searchService.Search(query, page, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	return marshalResult(map[string]interface{}{
		"data": results,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func handleGetRandomHadith(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hadith, err := hadithService.GetRandomHadith()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get random hadith: %v", err)), nil
	}
	return marshalResult(hadith)
}
