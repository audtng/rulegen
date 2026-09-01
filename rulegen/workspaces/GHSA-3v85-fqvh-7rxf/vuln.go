package main

}

func (echoService *EchoService) LikeEcho(ctx context.Context, id string) error {
	if err := echoService.transactor.Run(ctx, func(txCtx context.Context) error {
		return echoService.echoRepository.LikeEcho(txCtx, id)
	}); err != nil {
	if cleaned == "" {
		return nil, errors.New(commonModel.INVALID_PARAMS)
	}

	existing, err := echoService.echoRepository.GetTagsByNames(ctx, []string{cleaned})
	if err != nil {
	var names []string
	for _, tag := range echo.Tags {
		name := strings.TrimSpace(strings.TrimPrefix(tag.Name, "#"))
		if name != "" {
			names = append(names, name)
		}
	}

	existingTags, err := echoService.echoRepository.GetTagsByNames(ctx, names)
	}, nil
}

func normalizeEchoExtension(ext *model.EchoExtension) (*model.EchoExtension, error) {
	if ext == nil {
		return nil, nil
	"context"
	"errors"
	"fmt"
	"strings"
	"time"


				if len(msg.Tags) > 0 {
					for _, tag := range msg.Tags {
						renderedContent = fmt.Appendf(renderedContent, "<br /><span class=\"tag\">#%s</span>", tag.Name)
					}
				}

	doc := p.Parse(md)

	// 创建 HTML 渲染器
	htmlFlags := html.CommonFlags |
		html.Safelink |
		html.HrefTargetBlank |
		html.NoopenerLinks |
		html.NoreferrerLinks
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

