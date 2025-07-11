import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChatStore = defineStore('chat', () => {
  // --- State ---
  const newChatMessages = ref([]);
  const newGroupChatMessages = ref([]);
  const unreadMessages = ref([]);
  const unreadMsgsStatsFromDB = ref([]);
  const openChats = ref([]);
  const chatUserList = ref([]);

  const getMessages = (id) => (receiverId, type) => {
    if (type === 'PERSON') {
      return newChatMessages.value.filter(
        (e) =>
          (e.receiverId === receiverId && e.senderId === id) ||
          (e.receiverId === id && e.senderId === receiverId)
      );
    } else {
      return newGroupChatMessages.value.filter((msg) => msg.receiverId === receiverId);
    }
  };

  const getUnreadMessagesCount = (id) => (userId) => {
    return unreadMessages.value.filter(
      (msg) => msg.senderId === userId && msg.receiverId === id
    ).length;
  };

  const getUnreadGroupMessagesCount = (groupId) => {
    return unreadMessages.value.filter((msg) => msg.receiverId === groupId && msg.type === 'GROUP').length;
  };

  const getUnreadMsgsCountFromDB = (userId) => {
    if (!unreadMsgsStatsFromDB.value) return 0;
    const userMsgObj = unreadMsgsStatsFromDB.value.find((msg) => msg.id === userId);
    return userMsgObj ? userMsgObj.unreadMsgCount : 0;
  };

  async function fetchUnreadMessages() {
    try {
        const response = await fetch('http://localhost:8081/unreadMessages', {
            credentials: 'include'
        });
        const data = await response.json();
        if (data.type === 'Error') {
            unreadMsgsStatsFromDB.value = null;
        } else {
            unreadMsgsStatsFromDB.value = data.chatStats;
        }
    } catch (error) {
        console.error("Failed to fetch unread messages:", error);
        unreadMsgsStatsFromDB.value = null;
    }
  }

  function addNewChatMessage(payload) {
    if (payload.type === 'PERSON') {
      newChatMessages.value.push(payload);
    } else {
      newGroupChatMessages.value.push(payload);
    }
  }

  function addUnreadChatMessage(payload) {
    unreadMessages.value.push(payload);
  }

  function removeUnreadMessages(payload) {
    if (payload.type === 'GROUP') {
      unreadMessages.value = unreadMessages.value.filter((msg) => msg.receiverId !== payload.receiverId);
    } else {
      unreadMessages.value = unreadMessages.value.filter(
        (msg) => !(msg.type === 'PERSON' && msg.senderId === payload.receiverId)
      );
    }
  }

  function addNewChat(chatBox) {
    // Avoid adding duplicate chats
    if (!openChats.value.some(chat => chat.name === chatBox.name)) {
        openChats.value.push(chatBox);
    }
  }

  function removeChat(name) {
    openChats.value = openChats.value.filter((chat) => chat.name !== name);
  }

  function clearOpenChats() {
    openChats.value = [];
  }

  async function fetchChatUserList(userId) {
    try {
        const response = await fetch(`http://localhost:8081/chatList?userId=${userId}`, {
            credentials: 'include'
        });
        const data = await response.json();
        chatUserList.value = data.users;
    } catch (error) {
        console.error("Failed to fetch chat user list:", error);
        chatUserList.value = [];
    }
  }
  function resetChatWithUser(currentUserId, otherUserId) {
    newChatMessages.value = newChatMessages.value.filter(
      msg =>
        !(
          (msg.senderId === currentUserId && msg.receiverId === otherUserId) ||
          (msg.senderId === otherUserId && msg.receiverId === currentUserId)
        )
    );
  }
  return {
    newChatMessages,
    newGroupChatMessages,
    unreadMessages,
    unreadMsgsStatsFromDB,
    openChats,
    chatUserList,
    getMessages,
    getUnreadMessagesCount,
    getUnreadGroupMessagesCount,
    getUnreadMsgsCountFromDB,
    fetchUnreadMessages,
    addNewChatMessage,
    addUnreadChatMessage,
    removeUnreadMessages,
    addNewChat,
    removeChat,
    clearOpenChats,
    fetchChatUserList,
    resetChatWithUser
  };
});