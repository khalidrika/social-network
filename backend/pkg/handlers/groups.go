package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"social-network/pkg/models"
	"social-network/pkg/utils"
	ws "social-network/pkg/wsServer"
	"strings"
)

/* -------------------------------------------------------------------------- */
/*                                  helpers                                   */
/* -------------------------------------------------------------------------- */

// checkGroupAccess verifies if user can access a group based on its privacy settings
func (handler *Handler) checkGroupAccess(groupId, userId string) (bool, error) {
	// Get group data to check privacy
	group, err := handler.repos.GroupRepo.GetData(groupId)
	if err != nil {
		return false, err
	}

	// If group is public, anyone can access
	if group.Privacy == "public" {
		return true, nil
	}

	// If group is private, only members and admin can access
	if group.Privacy == "private" {
		isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
		if err != nil {
			return false, err
		}
		if isAdmin {
			return true, nil
		}

		isMember, err := handler.repos.GroupRepo.IsMember(groupId, userId)
		if err != nil {
			return false, err
		}
		return isMember, nil
	}

	return false, nil
}

/* -------------------------------------------------------------------------- */
/*                                  get data                                  */
/* -------------------------------------------------------------------------- */

func (handler *Handler) AllGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request all groups + relations
	groups, errGroups := handler.repos.GroupRepo.GetAllAndRelations(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns all groups that current user is a member of or admin
func (handler *Handler) UserGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	// request user Groups
	groups, errGroups := handler.repos.GroupRepo.GetUserGroups(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns all groups that specified user is a member of or admin
func (handler *Handler) OtherUserGroups(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access user id
	query := r.URL.Query()
	userId := query.Get("userId")
	if userId == "" { //check if user id provided in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// request user Groups
	groups, errGroups := handler.repos.GroupRepo.GetUserGroups(userId)
	if errGroups != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithGroups(w, groups, 200)
}

// returns info about group - > name, description, id and administrator id
// also includes group status for current user -> admin / member or pending member request
func (handler *Handler) GroupInfo(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}

	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)

	// Check if user can access this group
	canAccess, err := handler.checkGroupAccess(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !canAccess {
		utils.RespondWithError(w, "Access denied to private group", 403)
		return
	}

	group, err := handler.repos.GroupRepo.GetData(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* --------------- // add additional data about member status --------------- */
	// check if admin, or member or member request pending
	if userId == group.AdminID {
		group.Administrator = true
		group.Member = true // Admin is also a member
	} else {
		group.Member, err = handler.repos.GroupRepo.IsMember(group.ID, userId)
		if err != nil {
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
		if !group.Member {
			notification := models.Notification{
				TargetID: group.ID,
				Type:     "GROUP_REQUEST",
				Content:  userId,
			}
			group.RequestPending, err = handler.repos.NotifRepo.CheckIfExists(notification)
			if err != nil {
				utils.RespondWithError(w, "Error on getting data", 200)
				return
			}
		}
	}
	utils.RespondWithGroups(w, []models.Group{group}, 200)
}

// returns list of all group members and administrator
func (handler *Handler) GroupMembers(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}

	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)

	// Check if user can access this group
	canAccess, err := handler.checkGroupAccess(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !canAccess {
		utils.RespondWithError(w, "Access denied to private group", 403)
		return
	}

	members, err := handler.repos.GroupRepo.GetMembers(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithUsers(w, members, 200)
}

// returns list of events for group
// GET request with group_id included in query
// in case if user is not a member of the group, returns error
func (handler *Handler) GroupEvents(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on reading group ID", 200)
		return
	}
	// check if current user is a member or admin  of the group
	var isMember = false
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on reading role", 200)
		return
	}
	if !isAdmin {
		isMember, err = handler.repos.GroupRepo.IsMember(groupId, userId)
		if err != nil {
			utils.RespondWithError(w, "Error on checking if is group member", 200)
			return
		}
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	// current user is a member or admin -> get events with responses
	events, err := handler.repos.EventRepo.GetGroupEventsWithResponses(groupId, userId)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, "Error on getting event data", 200)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// returns all posts that belongs to group
func (handler *Handler) GroupPosts(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}

	// Check if user can access this group
	canAccess, err := handler.checkGroupAccess(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !canAccess {
		utils.RespondWithError(w, "Access denied to private group", 403)
		return
	}

	/* ------------- current user can access group -> get posts ------------- */
	posts, err := handler.repos.PostRepo.GetGroupPosts(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get post author info attached
	if err = AttachAuthors(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// Get comment info for each post
	if err = AttachComments(handler, &posts); err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	utils.RespondWithPosts(w, posts, 200)
}

// returns pending requests to join to group, only for admin
// for others respond with error
func (handler *Handler) GroupRequests(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	/* ------------------------- check if user is admin ------------------------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isAdmin {
		utils.RespondWithError(w, "Unauthorized", 200)
		return
	}
	/* ---------------------- get pending requests form db ---------------------- */
	notifications, err := handler.repos.NotifRepo.GetGroupRequests(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	for i := 0; i < len(notifications); i++ {
		notifications[i].User, err = handler.repos.UserRepo.GetDataMin(notifications[i].Sender)
		if err != nil {
			utils.RespondWithError(w, "Error on getting data", 200)
			return
		}
	}
	utils.RespondWithNotifications(w, notifications, 200)
}

func (handler *Handler) CancelGroupRequests(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	currentUserId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	// delete notification
	notif := models.Notification{
		Type:     "GROUP_REQUEST",
		TargetID: groupId,
		Content:  currentUserId,
	}
	if err := handler.repos.NotifRepo.DeleteByType(notif); err != nil {
		utils.RespondWithError(w, "Error on canceling request", 200)
		return
	}
	utils.RespondWithSuccess(w, "gROUP request canceled successfuly", 200)
}

/* -------------------------------------------------------------------------- */
/*                                save new data                               */
/* -------------------------------------------------------------------------- */

func (handler *Handler) NewGroup(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		utils.RespondWithError(w, "Error parsing form data", 200)
		return
	}

	// Create new group from form data
	var newGroup models.Group
	newGroup.Name = r.FormValue("name")
	newGroup.Description = r.FormValue("description")
	newGroup.Privacy = r.FormValue("privacy")
	
	// Handle image upload
	file, fileHeader, err := r.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()
		
		// Save the uploaded file
		imagePath, err := utils.SaveUploadedFile(file, fileHeader, "imageUpload")
		if err != nil {
			utils.RespondWithError(w, "Error saving image", 200)
			return
		}
		newGroup.Image = imagePath
	}

	// Parse invitations JSON if provided
	invitationsJSON := r.FormValue("invitations")
	if invitationsJSON != "" {
		err = json.Unmarshal([]byte(invitationsJSON), &newGroup.Invitations)
		if err != nil {
			utils.RespondWithError(w, "Error parsing invitations", 200)
			return
		}
	}

	/* ------------------------- attach addidtional data ------------------------ */
	// access user id
	newGroup.AdminID = r.Context().Value(utils.UserKey).(string)
	newGroup.ID = utils.UniqueId()

	/* ------------------------------- save in db ------------------------------- */
	err = handler.repos.GroupRepo.New(newGroup)
	if err != nil {
		utils.RespondWithError(w, "Error on saving group", 200)
		return
	}
	/* ------------------------- save invitations in db ------------------------- */
	for i := 0; i < len(newGroup.Invitations); i++ {
		// save each follower in db
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: newGroup.Invitations[i],
			Type:     "GROUP_INVITE",
			Content:  newGroup.ID,
			Sender:   newGroup.AdminID,
		}
		err = handler.repos.NotifRepo.Save(newNotif)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

		for client := range wsServer.Clients {
			if client.ID == newNotif.TargetID {
				client.SendNotification(newNotif)
			}
		}

		// save as a new member of group
		// if err = handler.repos.GroupRepo.SaveMember(newGroup.Invitations[i], newGroup.ID); err != nil {
		// 	utils.RespondWithError(w, "Internal server error", 200)
		// 	return
		// }
	}
	newGroup.Administrator = true
	//NOTIFY WEBSOCKET ABOUT NEW NOTIFICATION
	utils.RespondWithGroups(w, []models.Group{newGroup}, 200)
}

// NOT TESTED
func (handler *Handler) NewGroupPost(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	err := r.ParseMultipartForm(3145728) // 3MB
	if err != nil {
		utils.RespondWithError(w, "Error in form validation", 200)
		return
	}
	userId := r.Context().Value(utils.UserKey).(string)
	/* ------------------------ create new post instance ------------------------ */
	newPost := models.Post{
		ID:       utils.UniqueId(),
		Content:  r.PostFormValue("body"),
		GroupID:  r.PostFormValue("groupId"),
		AuthorID: userId,
	}
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(newPost.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	isMember, err := handler.repos.GroupRepo.IsMember(newPost.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	/* ------------------------ save image in filesystem ------------------------ */
	newPost.ImagePath = utils.SaveImage(r)
	/* -------------------------- save post in database ------------------------- */
	errDB := handler.repos.PostRepo.New(newPost)
	if errDB != nil {
		utils.RespondWithError(w, "Error on saving post", 200)
		return
	}
	utils.RespondWithSuccess(w, "New post created", 200)
}

// handle when new user wants to join the group
func (handler *Handler) NewGroupRequest(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// access current user id
	userId := r.Context().Value(utils.UserKey).(string)
	// get group id from request
	query := r.URL.Query()
	groupId := query.Get("groupId")
	if groupId == "" { //check if group id exists in request
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	
	/* --------------------------- check if not admin --------------------------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if isAdmin {
		utils.RespondWithError(w, "Invalid request", 200)
		return
	}
	/* -------------------- create new notification instance -------------------- */
	notification := models.Notification{
		ID:       utils.UniqueId(),
		TargetID: groupId,
		Type:     "GROUP_REQUEST",
		Content:  userId,
		Sender:   userId,
	}
	
	/* --------------------- check if request already exists -------------------- */
	exists, err := handler.repos.NotifRepo.CheckIfExists(notification)
	if err != nil {
		utils.RespondWithError(w, "Error on saving request", 200)
		return
	}
	if exists {
		utils.RespondWithSuccess(w, "Request already saved", 200)
		return
	}
	/* ------------------------- save notification in db ------------------------ */
	err = handler.repos.NotifRepo.Save(notification)
	if err != nil {
		utils.RespondWithError(w, "Error on saving request", 200)
		return
	}

	//SEND MESSAGE TO TO GROUP ADMIN IF IS ONLINE
	admin, err := handler.repos.GroupRepo.GetAdmin(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error on finding group admin", 200)
		return
	}
	
	for client := range wsServer.Clients {
		if client.ID == admin {
			client.SendNotification(notification)
		}
	}
	utils.RespondWithSuccess(w, "Request saved successfuly", 200)
}

// NOT TESTED
// handle response from group administrator for requests to join group
// waits for requestId and response -accept/decline
func (handler *Handler) ResponseGroupRequest(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on data submittion", 200)
		return
	}
	/* --------------------------- read incoming data --------------------------- */
	type Response struct {
		GroupID   string `json:"groupId"`
		RequestID string `json:"requestId"`
		Response  string `json:"response"`
	}
	var response Response
	
	// Read body first to see what we're getting
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.RespondWithError(w, "Error reading request body", 200)
		return
	}
	
	// Parse the JSON
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		utils.RespondWithError(w, "Error parsing request body", 400)
		return
	}
	
	/* ---------------------- check if all fields provided ---------------------- */
	if response.Response != "accept" && response.Response != "decline" {
		utils.RespondWithError(w, "Invalid response. Must be 'accept' or 'decline'", 400)
		return
	}
	if response.GroupID == "" || response.RequestID == "" {
		utils.RespondWithError(w, "Missing required fields: groupId and requestId", 400)
		return
	}
	/* ------------------- chack if curren user is group admin ------------------ */
	// access user id
	userId := r.Context().Value(utils.UserKey).(string)
	
	if userId == "" {
		utils.RespondWithError(w, "Unauthorized - not logged in", 401)
		return
	}
	
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(response.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error checking admin status", 200)
		return
	}
	if !isAdmin {
		utils.RespondWithError(w, "Unauthorized", 200)
		return
	}
	
	/* ----------------------------- handle response ---------------------------- */
	// if accepted -> save as new member
	if response.Response == "accept" {
		// get id of member that requests to join
		joinerId, err := handler.repos.NotifRepo.GetUserFromRequest(response.RequestID)
		if err != nil {
			utils.RespondWithError(w, "Error retrieving user from request", 500)
			return
		}

		// Check if user is already a member to prevent duplicates
		isMember, err := handler.repos.GroupRepo.IsMember(response.GroupID, joinerId)
		if err != nil {
			utils.RespondWithError(w, "Error checking membership status", 500)
			return
		}

		// save as a new member of group only if they are not a member yet
		if !isMember {
			if err = handler.repos.GroupRepo.SaveMember(joinerId, response.GroupID); err != nil {
				utils.RespondWithError(w, "Error adding member to group", 500)
				return
			}
		}

		// if joiner online, send updated group status
		for client := range wsServer.Clients {
			if client.ID == joinerId {
				client.SendGroupRequestAccept(response.GroupID)
			}
		}
	}
	
	// delete from pending notification table
	if err = handler.repos.NotifRepo.Delete(response.RequestID); err != nil {
		utils.RespondWithError(w, "Error removing notification", 500)
		return
	}
	
	utils.RespondWithSuccess(w, "Group request processed successfully", 200)
}

// NOT TESTED
// waits for POST request with groupId and Invitation list included
func (handler *Handler) NewGroupInvite(wsServer *ws.Server, w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new Group
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	userId := r.Context().Value(utils.UserKey).(string)
	/* -------- check if current user is a member or admin  of the group -------- */
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(group.ID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	isMember, err := handler.repos.GroupRepo.IsMember(group.ID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error on getting data", 200)
		return
	}
	if !isMember && !isAdmin {
		utils.RespondWithError(w, "Not a member", 200)
		return
	}
	for i := 0; i < len(group.Invitations); i++ {
		// save each invitation in db
		newNotif := models.Notification{
			ID:       utils.UniqueId(),
			TargetID: group.Invitations[i],
			Type:     "GROUP_INVITE",
			Content:  group.ID,
			Sender:   userId,
		}
		err = handler.repos.NotifRepo.Save(newNotif)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}

		// search if user has open ws connection and send notification
		for client := range wsServer.Clients {
			if client.ID == group.Invitations[i] {
				client.SendNotification(newNotif)
			}
		}
	}
	utils.RespondWithSuccess(w, "Invitations saved", 200)
}

// NOT TESTED
func (handler *Handler) ResponseInviteRequest(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	/* ---------------------------- read incoming data --------------------------- */
	// Try to decode the JSON request to a new response
	type Response struct {
		RequestID string `json:"requestId"`
		Response  string `json:"response"` // ACCEPT or DECLINE
	}
	var resp Response
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		utils.RespondWithError(w, "Error on form submittion", 200)
		return
	}
	// get group id from notification
	groupId, err := handler.repos.NotifRepo.GetGroupId(resp.RequestID)
	userId := r.Context().Value(utils.UserKey).(string)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	if strings.ToUpper(resp.Response) == "ACCEPT" {
		// Check if user is already a member to prevent duplicates
		isMember, err := handler.repos.GroupRepo.IsMember(groupId, userId)
		if err != nil {
			utils.RespondWithError(w, "Internal server error", 200)
			return
		}
		
		// Only add as member if not already a member
		if !isMember {
			err = handler.repos.GroupRepo.SaveMember(userId, groupId)
			if err != nil {
				utils.RespondWithError(w, "Internal server error", 200)
				return
			}
		}
	}
	/* ----------------------- delete pending notification ---------------------- */
	err = handler.repos.NotifRepo.Delete(resp.RequestID)
	if err != nil {
		utils.RespondWithError(w, "Internal server error", 200)
		return
	}
	// notify websocket about notification changes
	utils.RespondWithSuccess(w, "Response successful", 200)
}

// JoinPublicGroup allows users to join public groups directly without admin approval
func (handler *Handler) JoinPublicGroup(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}

	// Get current user ID
	userId := r.Context().Value(utils.UserKey).(string)

	// Parse request body
	type JoinRequest struct {
		GroupID string `json:"groupId"`
	}
	var joinReq JoinRequest
	err := json.NewDecoder(r.Body).Decode(&joinReq)
	if err != nil {
		utils.RespondWithError(w, "Invalid request body", 400)
		return
	}

	if joinReq.GroupID == "" {
		utils.RespondWithError(w, "Group ID is required", 400)
		return
	}

	// Check if group exists and is public
	group, err := handler.repos.GroupRepo.GetData(joinReq.GroupID)
	if err != nil {
		utils.RespondWithError(w, "Group not found", 404)
		return
	}

	if group.Privacy != "public" {
		utils.RespondWithError(w, "Can only join public groups directly", 403)
		return
	}

	// Check if user is already an admin
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(joinReq.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error checking admin status", 500)
		return
	}
	if isAdmin {
		utils.RespondWithError(w, "You are already an admin of this group", 400)
		return
	}

	// Check if user is already a member
	isMember, err := handler.repos.GroupRepo.IsMember(joinReq.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error checking membership status", 500)
		return
	}
	if isMember {
		utils.RespondWithError(w, "You are already a member of this group", 400)
		return
	}

	// Add user as a member
	err = handler.repos.GroupRepo.SaveMember(userId, joinReq.GroupID)
	if err != nil {
		utils.RespondWithError(w, "Error joining group", 500)
		return
	}

	utils.RespondWithSuccess(w, "Successfully joined the group", 200)
}

// LeaveGroup allows members to leave a group (but not admins)
func (handler *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}

	// Get current user ID
	userId := r.Context().Value(utils.UserKey).(string)

	// Parse request body
	type LeaveRequest struct {
		GroupID string `json:"groupId"`
	}
	var leaveReq LeaveRequest
	err := json.NewDecoder(r.Body).Decode(&leaveReq)
	if err != nil {
		utils.RespondWithError(w, "Invalid request body", 400)
		return
	}

	if leaveReq.GroupID == "" {
		utils.RespondWithError(w, "Group ID is required", 400)
		return
	}

	// Check if group exists
	_, err = handler.repos.GroupRepo.GetData(leaveReq.GroupID)
	if err != nil {
		utils.RespondWithError(w, "Group not found", 404)
		return
	}

	// Check if user is the admin - admins cannot leave their own groups
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(leaveReq.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error checking admin status", 500)
		return
	}
	if isAdmin {
		utils.RespondWithError(w, "Admins cannot leave their own groups", 403)
		return
	}

	// Check if user is a member
	isMember, err := handler.repos.GroupRepo.IsMember(leaveReq.GroupID, userId)
	if err != nil {
		utils.RespondWithError(w, "Error checking membership status", 500)
		return
	}
	if !isMember {
		utils.RespondWithError(w, "You are not a member of this group", 400)
		return
	}

	// Remove user from the group
	err = handler.repos.GroupRepo.RemoveMember(userId, leaveReq.GroupID)
	if err != nil {
		utils.RespondWithError(w, "Error leaving group", 500)
		return
	}

	utils.RespondWithSuccess(w, "Successfully left the group", 200)
}

// CancelGroupInvite allows a group member to cancel a pending invitation.
func (handler *Handler) CancelGroupInvite(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "POST" {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}

	// Get current user ID (the one canceling the invite)
	currentUserId := r.Context().Value(utils.UserKey).(string)

	// Parse request body
	type CancelRequest struct {
		GroupID  string `json:"groupId"`
		TargetID string `json:"targetId"` // The user being uninvited
	}
	var cancelReq CancelRequest
	err := json.NewDecoder(r.Body).Decode(&cancelReq)
	if err != nil {
		utils.RespondWithError(w, "Invalid request body", 400)
		return
	}

	if cancelReq.GroupID == "" || cancelReq.TargetID == "" {
		utils.RespondWithError(w, "Group ID and Target ID are required", 400)
		return
	}

	// To cancel an invite, you must be a member or admin of the group.
	isMember, err := handler.repos.GroupRepo.IsMember(cancelReq.GroupID, currentUserId)
	if err != nil {
		utils.RespondWithError(w, "Error checking membership status", 500)
		return
	}
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(cancelReq.GroupID, currentUserId)
	if err != nil {
		utils.RespondWithError(w, "Error checking admin status", 500)
		return
	}

	if !isMember && !isAdmin {
		utils.RespondWithError(w, "You must be a member to cancel invites.", 403)
		return
	}

	// Create a notification model to find and delete the specific invitation.
	// We assume anyone in the group can cancel any invite.
	// This requires a new repo method: DeleteGroupInvite(targetId, groupId)
	if err := handler.repos.NotifRepo.DeleteGroupInvite(cancelReq.TargetID, cancelReq.GroupID); err != nil {
		utils.RespondWithError(w, "Error on canceling invite: "+err.Error(), 500)
		return
	}

	utils.RespondWithSuccess(w, "Group invite canceled successfully", 200)
}

// Check existing invitations for a group
func (handler *Handler) CheckGroupInvitations(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	if r.Method != "GET" {
		utils.RespondWithError(w, "Method not allowed", 405)
		return
	}

	// Get group ID from query
	groupId := r.URL.Query().Get("groupId")
	if groupId == "" {
		utils.RespondWithError(w, "Group ID is required", 400)
		return
	}

	userId := r.Context().Value(utils.UserKey).(string)

	// Check if current user is admin or member of the group
	isAdmin, err := handler.repos.GroupRepo.IsAdmin(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error checking admin status", 200)
		return
	}
	
	isMember, err := handler.repos.GroupRepo.IsMember(groupId, userId)
	if err != nil {
		utils.RespondWithError(w, "Error checking membership", 200)
		return
	}
	
	if !isAdmin && !isMember {
		utils.RespondWithError(w, "Access denied", 403)
		return
	}

	// Get all pending GROUP_INVITE notifications for this group
	notifications, err := handler.repos.NotifRepo.GetGroupInvites(groupId)
	if err != nil {
		utils.RespondWithError(w, "Error fetching notifications", 200)
		return
	}

	var existingInvites []string
	for _, notif := range notifications {
		existingInvites = append(existingInvites, notif.TargetID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"invitedUsers": existingInvites,
	})
}
