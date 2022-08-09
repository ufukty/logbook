package request


type UserCreate struct {
	Request struct {
		EmailAddress string `json:"email_address"`
		Salt         string `json:"salt"`
		Password     string `json:"password"`
	}
	Response struct {
		Status       bool   `json:"status"`
		ErrorMessage string `json:"error_message"`
	}
}

type UserDelete struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
	}
	Response struct {
		Status       bool   `json:"status"`
		ErrorMessage string `json:"error_message"`
	}
}

type BookmarkCreate struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
		UserId             string `json:"user_id"`
		TaskId             string `json:"task_id"`
		DisplayName        string `json:"display_name"`
		RootBookmark       bool   `json:"root_bookmark"`
	}
	Response struct {
	}
}

type BookmarkList struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
		UserId             string `json:"user_id"`
	}
	Response struct {
	}
}

type ChronologicalViewPlacementArray struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
		UserId             string `json:"user_id"`
		RootTaskId         string `json:"root_task_id"`
		Offset             int    `json:"offset"`
		Limit              int    `json:"limit"`
	}
	Response struct {
	}
}

type TaskCreate struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
		UserId             string `json:"user_id"`
		Content            string `json:"content"`
		SuperTaskId        string `json:"super_task_id"`
		CurrentRevisionId  string `json:"current_revision_id"`
	}
	Response struct {
	}
}

type TaskUpdateContent struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
		UserId             string `json:"user_id"`
		NewContent         string `json:"new_content"`
		CurrentRevisionId  string `json:"current_revision_id"`
	}
	Response struct {
	}
}

type TaskReattach struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
		UserId             string `json:"user_id"`
		NewSuperTaskId     string `json:"new_super_task_id"`
		CurrentRevisionId  string `json:"current_revision_id"`
	}
	Response struct {
	}
}

type AuthorizationRequiredRequestTemplate struct {
	Request struct {
		AuthorizationToken string `json:"authorization_token"`
		UserId             string `json:"user_id"`
	}
	Response struct {
	}
}

func (parameters *AuthorizationRequiredRequestTemplate) Mainngn() {
	authorization_token := AuthorizationRequiredRequestTemplate{}
	authorization_token.Request.AuthorizationToken = ""
	authorization_token.Request.UserId = ""
}