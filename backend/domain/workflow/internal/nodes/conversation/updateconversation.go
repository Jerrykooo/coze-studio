package conversation

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	wf "github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/execute"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ternary"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type UpdateConversation struct {
}

func NewUpdateConversation(_ context.Context) *UpdateConversation {
	return &UpdateConversation{}
}

func (c *UpdateConversation) Update(ctx context.Context, in map[string]any) (map[string]any, error) {

	var (
		execCtx     = execute.GetExeCtx(ctx)
		env         = ternary.IFElse(execCtx.ExeCfg.Mode == vo.ExecuteModeRelease, vo.Online, vo.Draft)
		appID       = execCtx.ExeCfg.AppID
		agentID     = execCtx.ExeCfg.AgentID
		version     = execCtx.ExeCfg.Version
		connectorID = execCtx.ExeCfg.ConnectorID
		userID      = execCtx.ExeCfg.Operator
	)
	cName, ok := in["conversationName"]
	if !ok {
		return nil, vo.WrapError(errno.ErrInvalidParameter, errors.New("conversation name is required"))
	}

	conversationName := cName.(string)

	ncName, ok := in["newConversationName"]
	if !ok {
		return nil, vo.WrapError(errno.ErrInvalidParameter, errors.New("new conversationName name is required"))
	}

	newConversationName := ncName.(string)

	if agentID != nil {
		return nil, vo.WrapError(errno.ErrConversationNodesNotAvailable, fmt.Errorf("in the agent scenario, update conversation is not available"))
	}

	if appID == nil {
		return nil, vo.WrapError(errno.ErrConversationNodesNotAvailable, errors.New("conversation update node, app id is required"))
	}

	_, existed, err := wf.GetRepository().GetConversationTemplate(ctx, env, vo.GetConversationTemplatePolicy{
		AppID:   appID,
		Name:    ptr.Of(conversationName),
		Version: ptr.Of(version),
	})
	if err != nil {
		return nil, err
	}

	if existed {
		return nil, vo.WrapError(errno.ErrConversationNodeInvalidOperation, fmt.Errorf("only conversation created through nodes are allowed to be modified or deleted"))
	}

	conversation, existed, err := wf.GetRepository().GetDynamicConversationByName(ctx, env, *appID, connectorID, userID, conversationName)
	if err != nil {
		return nil, err
	}

	if !existed {
		return map[string]any{
			"conversationId": "0",
			"isSuccess":      false,
			"isExisted":      false,
		}, nil
	}

	ncConversation, existed, err := wf.GetRepository().GetDynamicConversationByName(ctx, env, *appID, connectorID, userID, newConversationName)
	if err != nil {
		return nil, err
	}

	if existed {
		return map[string]any{
			"conversationId": strconv.FormatInt(ncConversation.ConversationID, 10),
			"isSuccess":      false,
			"isExisted":      true,
		}, nil
	}

	err = wf.GetRepository().UpdateDynamicConversationNameByID(ctx, env, conversation.ID, newConversationName)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"conversationId": strconv.FormatInt(conversation.ConversationID, 10),
		"isSuccess":      true,
		"isExisted":      false,
	}, nil

}
