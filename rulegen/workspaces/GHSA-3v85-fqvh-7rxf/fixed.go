package main

}

func (echoService *EchoService) LikeEcho(ctx context.Context, id string) error {
	echo, err := echoService.echoRepository.GetEchosById(ctx, id)
	if err != nil {
		return err
	}
	if echo == nil {
		return errors.New(commonModel.ECHO_NOT_FOUND)
	}
	// 与 GetEchoById 的可见性规则保持一致：匿名调用方禁止点赞私密 echo，
	// 已认证非管理员同样禁止；管理员（含 MCP 路径）允许。
	if echo.Private {
		userID := viewer.MustFromContext(ctx).UserID()
		if userID == "" {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}
		user, err := echoService.commonService.CommonGetUserByUserId(ctx, userID)
		if err != nil {
			return err
		}
		if !user.IsAdmin {
			return errors.New(commonModel.NO_PERMISSION_DENIED)
		}
	}

	if err := echoService.transactor.Run(ctx, func(txCtx context.Context) error {
		return echoService.echoRepository.LikeEcho(txCtx, id)
	}); err != nil {
	if cleaned == "" {
		return nil, errors.New(commonModel.INVALID_PARAMS)
	}
	if !isSafeTagName(cleaned) {
		return nil, errors.New(commonModel.INVALID_PARAMS)
	}

	existing, err := echoService.echoRepository.GetTagsByNames(ctx, []string{cleaned})
	if err != nil {
	var names []string
	for _, tag := range echo.Tags {
		name := strings.TrimSpace(strings.TrimPrefix(tag.Name, "#"))
		if name == "" {
			continue
		}
		if !isSafeTagName(name) {
			return errors.New(commonModel.INVALID_PARAMS)
		}
		names = append(names, name)
	}

	existingTags, err := echoService.echoRepository.GetTagsByNames(ctx, names)
	}, nil
}

// isSafeTagName 拒绝包含 HTML 元字符的标签名，配合 RSS 渲染端的 HTML 转义形成纵深防御
// （GHSA-3v85-fqvh-7rxf）。即使后续新增其他出口忘记转义，含 <>"'& 的标签也无法落库。
func isSafeTagName(name string) bool {
	return !strings.ContainsAny(name, "<>\"'&")
}

func normalizeEchoExtension(ext *model.EchoExtension) (*model.EchoExtension, error) {
	if ext == nil {
		return nil, nil
	"context"
	"errors"
	"fmt"
	stdhtml "html"
	"strings"
	"time"


				if len(msg.Tags) > 0 {
					for _, tag := range msg.Tags {
						// 标签名进入 RSS Atom <summary type="html"> 后会被订阅器二次解码并渲染成 HTML，
						// 必须先做 HTML 实体转义阻断 stored XSS（GHSA-3v85-fqvh-7rxf）。
						renderedContent = fmt.Appendf(
							renderedContent,
							"<br /><span class=\"tag\">#%s</span>",
							stdhtml.EscapeString(tag.Name),
						)
					}
				}

	doc := p.Parse(md)

	// 创建 HTML 渲染器
	// SkipHTML 丢弃 markdown 中的原始 HTML 块/内联，阻止 <script> 等标签被原样输出。
	// 该函数仅服务于 RSS Atom <summary type="html"> 渲染，前端正文走客户端 markdown-it，
	// 因此关闭原始 HTML 透传不会影响 Web UI。
	htmlFlags := html.CommonFlags |
		html.Safelink |
		html.HrefTargetBlank |
		html.NoopenerLinks |
		html.NoreferrerLinks |
		html.SkipHTML
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

