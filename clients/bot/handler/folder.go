package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
)

func (h *Handler) NewFolderHandler(ctx context.Context, mu *tgb.MessageUpdate) error {

	str := strings.Replace(mu.Text, "/folder", "", 1)

	if str == "" {
		return errors.New("empty command provided")
	}

	cloud, err := h.cloud.CreateFolder(str)

	if err != nil {
		return mu.Answer(errorAnswer(err.Error())).ParseMode(tg.HTML).DoVoid(ctx)
	}

	answer := tg.HTML.Text(
		tg.HTML.Line(
			"✨ Shiny folder",
			tg.HTML.Link(cloud.Title, cloud.PrivateURL),
			"has been created!",
		),
	)

	return mu.Answer(answer).
		ParseMode(tg.HTML).
		DoVoid(ctx)
}
