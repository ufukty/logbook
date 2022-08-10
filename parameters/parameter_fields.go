package parameters

import "time"

type UserCreate struct {
	Request struct {
		EmailAddress NonEmptyString `json:"email_address"`
		Salt         NonEmptyString `json:"salt"`
		Password     NonEmptyString `json:"password"`
	}
	Response struct {
		UserId UserId `json:"user_id"`
	}
}

type UserDelete struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
	}
	Response struct{}
}

type BookmarkCreate struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
		TaskId             TaskId         `json:"task_id"`
		DisplayName        NonEmptyString `json:"display_name"`
		RootBookmark       bool           `json:"root_bookmark"`
	}
	Response struct{}
}

type BookmarkList struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
	}
	Response struct {
		Bookmarks []struct {
			DisplayName NonEmptyString `json:"display_name"`
			TaskId      TaskId         `json:"task_id"`
		} `json:"bookmarks"`
	}
}

type PlacementArrayHierarchical struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
		RootTaskId         TaskId         `json:"root_task_id"`
		Offset             int            `json:"offset"`
		Limit              int            `json:"limit"`
	}
	Response struct {
		TotalItem int      `json:"total_item"`
		Offset    int      `json:"offset"`
		Limit     int      `json:"limit"`
		Items     []TaskId `json:"items"`
	}
}

type PlacementArrayChronological struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
		StartingTime       time.Time      `json:"root_task_id"`
		Offset             int            `json:"offset"`
		Limit              int            `json:"limit"`
	}
}

type TaskCreate struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
		Content            NonEmptyString `json:"content"`
		SuperTaskId        TaskId         `json:"super_task_id"`
		CurrentRevisionId  NonEmptyString `json:"current_revision_id"`
	}
	Response struct {
		TaskId TaskId `json:"task_id"`
	}
}

type TaskUpdateContent struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
		NewContent         NonEmptyString `json:"new_content"`
		CurrentRevisionId  NonEmptyString `json:"current_revision_id"`
	}
	Response struct{}
}

type TaskReattach struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
		NewSuperTaskId     TaskId         `json:"new_super_task_id"`
		CurrentRevisionId  NonEmptyString `json:"current_revision_id"`
	}
	Response struct{}
}

type TaskCreateByWrapping struct {
	Request struct {
		AuthorizationToken NonEmptyString
		UserId             UserId
		TaskIds            []TaskId
	}
	Response struct{}
}

type TemplateForAuthorizationRequiredRequests struct {
	Request struct {
		AuthorizationToken NonEmptyString `json:"authorization_token"`
		UserId             UserId         `json:"user_id"`
	}
	Response struct{}
}
