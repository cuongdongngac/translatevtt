package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	formatFlag := flag.String("format", "vtt", "Subtitle format to translate (vtt or srt)")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	log.Printf("Working directory: %s", cwd)

	// Read api.txt
	apiFile := filepath.Join(cwd, "api.txt")
	apiKey := ""
	modelName := "google/gemini-flash-1.5" // OpenRouter default model
	if b, err := os.ReadFile(apiFile); err == nil {
		lines := strings.Split(strings.TrimSpace(string(b)), "\n")
		if len(lines) > 0 {
			apiKey = strings.TrimSpace(lines[0])
		}
		if len(lines) > 1 {
			customModel := strings.TrimSpace(lines[1])
			if customModel != "" {
				modelName = customModel
			}
		}
	} else {
		apiKey = os.Getenv("OPENROUTER_API_KEY")
	}

	if apiKey == "" {
		log.Fatal("API Key not found in api.txt or OPENROUTER_API_KEY environment variable")
	}

	// Đọc instruction.txt
	instructionFile := filepath.Join(cwd, "instruction.txt")
	instructionText := ""
	if b, err := os.ReadFile(instructionFile); err == nil {
		instructionText = strings.TrimSpace(string(b))
	}

	ctx := context.Background()

	// Walk
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		extEn := ".en." + *formatFlag
		extVi := ".vi." + *formatFlag
		if !info.IsDir() && strings.HasSuffix(path, extEn) {
			viPath := strings.TrimSuffix(path, extEn) + extVi
			
			needTranslate := true
			if _, err := os.Stat(viPath); err == nil {
				// Tệp đã tồn tại, kiểm tra xem nội dung có phải tiếng Việt không
				content, err := os.ReadFile(viPath)
				if err == nil && isVietnamese(string(content)) {
					log.Printf("Đã tồn tại và là tiếng Việt %s, bỏ qua.", viPath)
					needTranslate = false
				} else {
					log.Printf("Tệp %s đã tồn tại nhưng không phải tiếng Việt. Tiến hành dịch lại...", viPath)
				}
			}
			
			if needTranslate {
				log.Printf("Translating: %s (using model: %s)", path, modelName)
				translateSubtitle(ctx, apiKey, modelName, path, viPath, instructionText, *formatFlag)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
	log.Println("Finished!")
}

func translateSubtitle(ctx context.Context, apiKey, modelName, enPath, viPath, extraInstruction, format string) {
	b, err := os.ReadFile(enPath)
	if err != nil {
		log.Printf("Error reading file %s: %v", enPath, err)
		return
	}
	content := string(b)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	blocks := strings.Split(content, "\n\n")

	var systemInstruction string
	if format == "srt" {
		systemInstruction = "You are a professional subtitle translator. Translate the following SubRip (SRT) file from English to Vietnamese.\n" +
			"ONLY translate the subtitle text. Keep the sequence numbers, timestamps, and empty lines exactly as they are.\n" +
			"Do NOT wrap your response in markdown formatting. Return raw valid SRT."
	} else {
		systemInstruction = "You are a professional subtitle translator. Translate the following WebVTT file from English to Vietnamese.\n" +
			"ONLY translate the subtitle text. Keep the `WEBVTT` header, timestamps, and empty lines exactly as they are.\n" +
			"Do NOT wrap your response in markdown formatting. Return raw valid WebVTT."
	}

	if extraInstruction != "" {
		systemInstruction += "\n\nAdditional instructions from user:\n" + extraInstruction
	}

	chunkSize := 40
	var finalTranslatedChunks []string

	for i := 0; i < len(blocks); i += chunkSize {
		end := i + chunkSize
		if end > len(blocks) {
			end = len(blocks)
		}
		
		// Bỏ qua các block trống
		isEmpty := true
		for _, b := range blocks[i:end] {
			if strings.TrimSpace(b) != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		chunkBlocks := blocks[i:end]
		chunkText := strings.Join(chunkBlocks, "\n\n")
		
		log.Printf("Translating chunk %d to %d (out of %d blocks)...", i+1, end, len(blocks))
		
		translatedChunk := callOpenRouter(ctx, apiKey, modelName, systemInstruction, chunkText)
		if translatedChunk == "" {
			log.Printf("Failed to translate chunk %d to %d. Skipping file %s", i+1, end, enPath)
			return
		}

		// Xóa thẻ WEBVTT dư thừa ở các chunk sau
		if i > 0 {
			translatedChunk = strings.TrimPrefix(translatedChunk, "WEBVTT\n")
			translatedChunk = strings.TrimPrefix(translatedChunk, "WEBVTT")
			translatedChunk = strings.TrimSpace(translatedChunk)
		}

		finalTranslatedChunks = append(finalTranslatedChunks, translatedChunk)
	}

	finalText := strings.Join(finalTranslatedChunks, "\n\n")

	err = os.WriteFile(viPath, []byte(finalText), 0644)
	if err != nil {
		log.Printf("Error saving file %s: %v", viPath, err)
	} else {
		log.Printf("Saved: %s", viPath)
	}
}

func callOpenRouter(ctx context.Context, apiKey, modelName, systemInstruction, content string) string {
	reqBody := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "system", "content": systemInstruction},
			{"role": "user", "content": content},
		},
	}
	
	jsonData, _ := json.Marshal(reqBody)
	maxRetries := 15
	var translatedText string

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return ""
		}
		
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("HTTP-Referer", "https://github.com/vtt-translator")
		req.Header.Set("X-Title", "VTT Translator")

		client := &http.Client{Timeout: 5 * time.Minute}
		resp, err := client.Do(req)
		
		if err != nil {
			log.Printf("Network error: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 429 {
			log.Printf("API Rate limit (429). Waiting 20 seconds... (Attempt %d/%d)", i+1, maxRetries)
			time.Sleep(20 * time.Second)
			continue
		} else if resp.StatusCode == 401 || resp.StatusCode == 403 {
			log.Fatalf("\n[FATAL ERROR] Invalid API Key or access denied. Check api.txt.\nDetails: %s", string(respBody))
		} else if resp.StatusCode != 200 {
			log.Printf("API Error %d: %s. Retrying in 15s... (Attempt %d/%d)", resp.StatusCode, string(respBody), i+1, maxRetries)
			time.Sleep(15 * time.Second)
			continue
		}

		var openRouterResp struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(respBody, &openRouterResp); err != nil {
			log.Printf("JSON parsing error: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if len(openRouterResp.Choices) > 0 {
			translatedText = openRouterResp.Choices[0].Message.Content
			break
		} else {
			log.Printf("Empty response from OpenRouter: %s", string(respBody))
			time.Sleep(10 * time.Second)
			continue
		}
	}

	if translatedText != "" {
		translatedText = strings.TrimSpace(translatedText)
		translatedText = strings.TrimPrefix(translatedText, "```vtt\n")
		translatedText = strings.TrimPrefix(translatedText, "```srt\n")
		translatedText = strings.TrimPrefix(translatedText, "```\n")
		translatedText = strings.TrimSuffix(translatedText, "\n```")
		translatedText = strings.TrimSpace(translatedText)
	}

	return translatedText
}

// Kiểm tra xem văn bản có chứa các ký tự đặc trưng của tiếng Việt hay không
func isVietnamese(text string) bool {
	viChars := "àáãảạăằắẵẳặâầấẫẩậđèéẽẻẹêềếễểệìíĩỉịòóõỏọôồốỗổộơờớỡởợùúũủụưừứữửựỳýỹỷỵ"
	lowerText := strings.ToLower(text)
	count := 0
	for _, char := range lowerText {
		if strings.ContainsRune(viChars, char) {
			count++
			// Nếu tìm thấy hơn 5 ký tự đặc trưng, khả năng rất cao là tiếng Việt
			if count > 5 {
				return true
			}
		}
	}
	return false
}
