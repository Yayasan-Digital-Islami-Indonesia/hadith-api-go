package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
)

func handleGetBooks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	books, err := bookService.GetAllBooks()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get books: %v", err)), nil
	}

	data, err := json.Marshal(books)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal books: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func handleGetBook(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	identifier, ok := request.Params.Arguments["identifier"].(string)
	if !ok {
		return mcp.NewToolResultError("identifier must be a string"), nil
	}

	book, err := bookService.GetBookOrBySlug(identifier)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get book: %v", err)), nil
	}
	if book == nil {
		return mcp.NewToolResultError("Book not found"), nil
	}

	data, err := json.Marshal(book)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal book: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func handleGetChapters(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	bookIDStr, ok := request.Params.Arguments["book_id"].(string)
	if !ok {
		return mcp.NewToolResultError("book_id must be a string"), nil
	}

	bookID, err := strconv.ParseUint(bookIDStr, 10, 32)
	if err != nil {
		return mcp.NewToolResultError("Invalid book_id"), nil
	}

	chapters, err := chapterService.GetChaptersByBook(uint(bookID))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get chapters: %v", err)), nil
	}

	data, err := json.Marshal(chapters)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal chapters: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func handleGetChapterHadiths(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	chapterID, ok := request.Params.Arguments["chapter_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("chapter_id must be a number"), nil
	}

	page := 1
	if p, ok := request.Params.Arguments["page"].(float64); ok && p > 0 {
		page = int(p)
	}

	limit := 20
	if l, ok := request.Params.Arguments["limit"].(float64); ok && l > 0 && l <= 100 {
		limit = int(l)
	}

	hadiths, total, err := hadithService.GetHadithsByChapter(uint(chapterID), page, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get hadiths: %v", err)), nil
	}

	result := map[string]interface{}{
		"data": hadiths,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func handleGetHadith(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.Params.Arguments["id"].(float64)
	if !ok {
		return mcp.NewToolResultError("id must be a number"), nil
	}

	hadith, err := hadithService.GetHadith(uint(id))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get hadith: %v", err)), nil
	}
	if hadith == nil {
		return mcp.NewToolResultError("Hadith not found"), nil
	}

	data, err := json.Marshal(hadith)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal hadith: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func handleGetHadithByNumber(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	bookID, ok := request.Params.Arguments["book_id"].(float64)
	if !ok {
		return mcp.NewToolResultError("book_id must be a number"), nil
	}

	number, ok := request.Params.Arguments["number"].(float64)
	if !ok {
		return mcp.NewToolResultError("number must be a number"), nil
	}

	hadith, err := hadithService.GetHadithByBookAndNumber(uint(bookID), int(number))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get hadith: %v", err)), nil
	}
	if hadith == nil {
		return mcp.NewToolResultError("Hadith not found"), nil
	}

	data, err := json.Marshal(hadith)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal hadith: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func handleSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return mcp.NewToolResultError("query must be a non-empty string"), nil
	}

	page := 1
	if p, ok := request.Params.Arguments["page"].(float64); ok && p > 0 {
		page = int(p)
	}

	limit := 20
	if l, ok := request.Params.Arguments["limit"].(float64); ok && l > 0 && l <= 100 {
		limit = int(l)
	}

	results, total, err := searchService.Search(query, page, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	result := map[string]interface{}{
		"data": results,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func handleGetRandomHadith(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hadith, err := hadithService.GetRandomHadith()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get random hadith: %v", err)), nil
	}

	data, err := json.Marshal(hadith)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal hadith: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}
