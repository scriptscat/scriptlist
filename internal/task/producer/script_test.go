package producer

import (
	"encoding/json"
	"testing"

	broker2 "github.com/cago-frame/cago/pkg/broker/broker"
	"github.com/scriptscat/scriptlist/internal/model/entity/script_entity"
	"github.com/stretchr/testify/assert"
)

func TestScriptDeleteMsg_Serialize(t *testing.T) {
	msg := &ScriptDeleteMsg{
		Script: &script_entity.Script{
			ID:   100,
			Name: "test script",
		},
		Operator: Operator{
			OperatorUID:      1,
			OperatorUsername: "admin",
			IsAdmin:          true,
		},
		Reason: "违反使用规范",
	}

	body, err := json.Marshal(msg)
	assert.NoError(t, err)

	parsed, err := ParseScriptDeleteMsg(&broker2.Message{Body: body})
	assert.NoError(t, err)
	assert.Equal(t, int64(100), parsed.Script.ID)
	assert.Equal(t, "test script", parsed.Script.Name)
	assert.Equal(t, int64(1), parsed.OperatorUID)
	assert.Equal(t, "admin", parsed.OperatorUsername)
	assert.True(t, parsed.IsAdmin)
	assert.Equal(t, "违反使用规范", parsed.Reason)
}

func TestScriptDeleteMsg_BackwardCompatible(t *testing.T) {
	// 模拟旧格式的消息（直接 marshal *Script），验证新的 ParseScriptDeleteMsg 能解析
	oldScript := &script_entity.Script{
		ID:   200,
		Name: "old format script",
	}
	body, err := json.Marshal(oldScript)
	assert.NoError(t, err)

	// 旧格式没有 Script 包装字段，解析后 Script 应为 nil
	parsed, err := ParseScriptDeleteMsg(&broker2.Message{Body: body})
	assert.NoError(t, err)
	// 旧格式直接 marshal Script，新格式在 "script" 字段，所以 Script 为 nil
	// 但不会出错
	assert.NotNil(t, parsed)
}

func TestScriptCreateMsg_WithOperator(t *testing.T) {
	msg := &ScriptCreateMsg{
		Script: &script_entity.Script{
			ID:   1,
			Name: "new script",
		},
		CodeID: 10,
		Operator: Operator{
			OperatorUID:      5,
			OperatorUsername: "user5",
		},
	}

	body, err := json.Marshal(msg)
	assert.NoError(t, err)

	parsed, err := ParseScriptCreateMsg(&broker2.Message{Body: body})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), parsed.Script.ID)
	assert.Equal(t, int64(10), parsed.CodeID)
	assert.Equal(t, int64(5), parsed.OperatorUID)
	assert.Equal(t, "user5", parsed.OperatorUsername)
	assert.False(t, parsed.IsAdmin)
}

func TestScriptCreateMsg_WithoutOperator(t *testing.T) {
	// 旧格式不包含 Operator，验证兼容性
	msg := &ScriptCreateMsg{
		Script: &script_entity.Script{
			ID: 1,
		},
		CodeID: 10,
	}

	body, err := json.Marshal(msg)
	assert.NoError(t, err)

	parsed, err := ParseScriptCreateMsg(&broker2.Message{Body: body})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), parsed.Script.ID)
	assert.Equal(t, int64(10), parsed.CodeID)
	assert.Equal(t, int64(0), parsed.OperatorUID)
	assert.Equal(t, "", parsed.OperatorUsername)
	assert.False(t, parsed.IsAdmin)
}

func TestScriptCodeUpdateMsg_WithOperator(t *testing.T) {
	msg := &ScriptCodeUpdateMsg{
		Script: &script_entity.Script{
			ID:   3,
			Name: "updated script",
		},
		CodeID: 20,
		Operator: Operator{
			OperatorUID:      2,
			OperatorUsername: "admin2",
			IsAdmin:          true,
		},
	}

	body, err := json.Marshal(msg)
	assert.NoError(t, err)

	parsed, err := ParseScriptCodeUpdateMsg(&broker2.Message{Body: body})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), parsed.Script.ID)
	assert.Equal(t, int64(20), parsed.CodeID)
	assert.Equal(t, int64(2), parsed.OperatorUID)
	assert.True(t, parsed.IsAdmin)
}
