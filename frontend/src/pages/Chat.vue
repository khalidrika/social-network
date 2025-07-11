<template>
  <div class="chat-container">
    <div class="chat-layout">
      <div class="friends-list-section">
        <div class="header">
          <div class="header-icon-wrapper">
            <svg class="header-icon" fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"></path></svg>
          </div>
          <h1 class="header-title">Discussions</h1>
        </div>

        <div v-if="currentUser" class="profile-card">
          <div class="profile-avatar-wrapper">
            <img :src="getFullImageUrl(currentUser.avatar || currentUser.imagePath)" alt="Avatar" class="profile-avatar" />
          </div>
          <div class="profile-name">{{ currentUser.nickname || `${currentUser.firstName} ${currentUser.lastName}` }}</div>
          <div class="profile-status">En ligne</div>
        </div>

        <div class="friends-list-container">
          <div class="friends-list-header">
            <span class="friends-list-title">Amis</span>
            <span class="friends-list-count">{{ following.length }}</span>
          </div>
          <div class="friends-list custom-scrollbar">
            <button
              v-for="friend in following"
              :key="friend.id"
              @click="selectConversation(friend, 'PERSON')"
              class="friend-item"
              :class="{ 'active': selectedConversation?.id === friend.id && selectedConversation.type === 'PERSON' }"
            >
              <div class="friend-avatar-wrapper">
                <img :src="getFullImageUrl(friend.avatar || friend.imagePath)" alt="Avatar" class="friend-avatar" />
                <span v-if="getUnreadCountForFriend(friend.id) > 0" class="unread-badge">{{ getUnreadCountForFriend(friend.id) }}</span>
              </div>
              <div class="friend-details">
                <div class="friend-name">{{ friend.nickname || friend.firstName }}</div>
              </div>
            </button>
          </div>
        </div>

        <div class="friends-list-container">
          <div class="friends-list-header">
            <span class="friends-list-title">Groupes</span>
            <span class="friends-list-count">{{ groupStore.userGroups.length }}</span>
          </div>
          <div class="friends-list custom-scrollbar">
            <button
              v-for="group in groupStore.userGroups"
              :key="group.id"
              @click="selectConversation(group, 'GROUP')"
              class="friend-item"
              :class="{ 'active': selectedConversation?.id === group.id && selectedConversation.type === 'GROUP' }"
            >
              <div class="friend-avatar-wrapper">
                <img :src="getFullImageUrl(group.image)" alt="Avatar" class="friend-avatar" />
                 <span v-if="chatStore.getUnreadGroupMessagesCount(group.id) > 0" class="unread-badge">{{ chatStore.getUnreadGroupMessagesCount(group.id) }}</span>
              </div>
              <div class="friend-details">
                <div class="friend-name">{{ group.name }}</div>
              </div>
            </button>
          </div>
        </div>
      </div>

      <div class="chat-area-wrapper">
        <div class="chat-box">
          <div v-if="selectedConversation" class="active-chat">
            <div class="chat-header">
              <div class="chat-header-user">
                 <img :src="getFullImageUrl(selectedConversation.avatar || selectedConversation.imagePath || selectedConversation.image)" alt="Avatar" class="chat-header-avatar" />
                <div class="chat-header-text">
                  <span class="chat-header-name">{{ selectedConversation.name || selectedConversation.nickname || selectedConversation.firstName }}</span>
                  <span class="chat-header-status" v-if="selectedConversation.type === 'PERSON'">En ligne</span>
                   <span class="chat-header-status" v-else>{{ selectedConversation.memberCount }} Membres</span>
                </div>
              </div>
            </div>

            <div ref="messageContainer" class="messages-container custom-scrollbar">
              <div class="messages-list">
                <div class="messages-grid">
                  <template v-for="(message, index) in chatMessages" :key="message.id || index">
                    <div v-if="message.senderId !== currentUser.id" class="message-wrapper received">
                      <div class="message-layout">
                        <div class="message-avatar-wrapper">
                           <img :src="getFullImageUrl(message.sender?.avatar || message.sender?.imagePath)" alt="Avatar" class="message-avatar" />
                        </div>
                        <div class="message-bubble received"><div>{{ message.content || message.text }}</div></div>
                      </div>
                    </div>
                    <div v-else class="message-wrapper sent">
                      <div class="message-layout sent">
                        <div class="message-bubble sent"><div>{{ message.content || message.text }}</div></div>
                         <div class="message-avatar-wrapper">
                           <img :src="getFullImageUrl(currentUser.avatar || currentUser.imagePath)" alt="Avatar" class="message-avatar" />
                        </div>
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </div>

            <div class="message-input-section">
                <div class="emojis-container">
                    <button @click="toggleEmojis" class="emoji-toggle-btn">😊</button>
                    <div v-if="showEmojis" class="emojis-picker">
                        <span v-for="emoji in emojiList" :key="emoji" @click="addEmoji(emoji)" class="emoji-item">{{ emoji }}</span>
                    </div>
                </div>
                <input
                    type="text"
                    v-model="newMessage"
                    @keyup.enter="sendMessage"
                    class="message-input"
                    placeholder="Écrivez votre message..."
                />
                <button @click="sendMessage" class="send-button">
                    <span class="send-button-icon-wrapper">
                        <svg class="send-button-icon" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path d="M10.894 2.553a1 1 0 00-1.788 0l-7 14a1 1 0 001.169 1.409l5-1.429A1 1 0 009 15.571V11a1 1 0 112 0v4.571a1 1 0 00.725.962l5 1.428a1 1 0 001.17-1.408l-7-14z"></path></svg>
                    </span>
                </button>
            </div>
          </div>
          <div v-else class="empty-chat-container">
            <div class="empty-chat-icon-wrapper">
              <svg class="empty-chat-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
            </div>
            <h2 class="empty-chat-title">Sélectionnez une conversation</h2>
            <p class="empty-chat-text">Choisissez un ami ou un groupe dans la liste pour commencer à discuter.</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { useUserStore } from '../stores/userStore';
import { useFollowStore } from '../stores/followStore';
import { useChatStore } from '../stores/chatStore';
import { useGroupStore } from '../stores/groupStore';
import { storeToRefs } from 'pinia';
import wsService from '../services/websocketService.js';

const userStore = useUserStore();
const followStore = useFollowStore();
const chatStore = useChatStore();
const groupStore = useGroupStore();

const { user: currentUser } = storeToRefs(userStore);
const { following } = storeToRefs(followStore);

const selectedConversation = ref(null);
const newMessage = ref('');
const messageContainer = ref(null);
const showEmojis = ref(false);
const removeWsListener = ref(null);
const historicalMessages = ref([]);

const chatMessages = computed(() => {
  if (!selectedConversation.value || !currentUser.value) return [];
  const { id, type } = selectedConversation.value;
  const newMsgs = chatStore.getMessages(currentUser.value.id)(id, type);
  return [...historicalMessages.value, ...newMsgs];
});

const emojiList = ['😊', '😂', '❤️', '👍', '🤔', '🎉', '😢', '😠'];

function getFullImageUrl(path) {
  if (!path) return 'https://placehold.co/120x120/e879c6/white?text=User';
  if (path.startsWith('http')) return path;
  try {
    return new URL(path.replace(/^\/+/, ''), 'http://localhost:8081').href;
  } catch (e) {
    console.error("Invalid image path:", path);
    return 'https://placehold.co/120x120/e879c6/white?text=Error';
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight;
    }
  });
};

const fetchPreviousMessages = async (conversation) => {
  if (!conversation) {
    historicalMessages.value = [];
    return;
  }
  try {
    const response = await fetch("http://localhost:8081/messages", {
        method: "POST",
        credentials: "include",
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            receiverId: conversation.id,
            type: conversation.type
        })
    });
    const data = await response.json();
    historicalMessages.value = data.chatMessage || [];
  } catch (error) {
    console.error("Failed to fetch previous messages:", error);
    historicalMessages.value = [];
  }
};

const selectConversation = async (conversation, type) => {
  selectedConversation.value = { ...conversation, type };
  
  chatStore.removeUnreadMessages({ receiverId: conversation.id, type });

  if (currentUser.value) {
    const currentUserId = currentUser.value.id;
    if (type === 'PERSON') {
      chatStore.newChatMessages = chatStore.newChatMessages.filter(msg => 
          !((msg.senderId === currentUserId && msg.receiverId === conversation.id) || 
            (msg.senderId === conversation.id && msg.receiverId === currentUserId))
      );
    } else {
      chatStore.newGroupChatMessages = chatStore.newGroupChatMessages.filter(msg => msg.receiverId !== conversation.id);
    }
  }
  
  await fetchPreviousMessages(selectedConversation.value); 
  
  scrollToBottom();
};

const handleIncomingMessage = (data) => {
  if (data && data.action === 'chat' && data.chatMessage) {
    const message = data.chatMessage;
    chatStore.addNewChatMessage(message);
    if (!selectedConversation.value || selectedConversation.value.id !== message.receiverId) {
      chatStore.addUnreadChatMessage(message);
    } else {
      scrollToBottom();
    }
  } else if (data && data.action === 'notification' && data.notification) {
        const message = data.notification;
        chatStore.addNewChatMessage(message);
        if (selectedConversation.value && (message.senderId === selectedConversation.value.id || message.receiverId === selectedConversation.value.id)) {
            scrollToBottom();
        }
    }
};

const sendMessage = async () => {
  const messageContent = newMessage.value.trim();
  if (messageContent === '' || !selectedConversation.value || !currentUser.value) {
    return;
  }
  const msgObj = {
    action: 'chat',
    receiverId: selectedConversation.value.id,
    content: messageContent,
    type: selectedConversation.value.type,
    senderId: currentUser.value.id,
    sender: { 
        id: currentUser.value.id,
        nickname: currentUser.value.nickname,
        avatar: currentUser.value.avatar,
        imagePath: currentUser.value.imagePath,
    }
  };
  try {
    const response = await fetch("http://localhost:8081/newMessage", {
        body: JSON.stringify(msgObj),
        method: "POST",
        credentials: "include",
        headers: { 'Content-Type': 'application/json' }
    });
    const data = await response.json();

    if (data.type === "Success") {
        // chatStore.addNewChatMessage({ ...msgObj, sender: { nickname: currentUser.value.nickname } }); // REMOVE THIS LINE
        newMessage.value = '';
        showEmojis.value = false;
        scrollToBottom();
    } else {
        console.error("Failed to send message:", data.message);
    }
  } catch (error) {
    console.error("Error sending message:", error);
  }
  wsService.send(msgObj);
};

const toggleEmojis = () => {
  showEmojis.value = !showEmojis.value;
};

const addEmoji = (emoji) => {
  newMessage.value += emoji;
};

const getUnreadCountForFriend = (friendId) => {
  if (!currentUser.value) return 0;
  return chatStore.getUnreadMessagesCount(currentUser.value.id)(friendId);
};

onMounted(async () => {
  try {
    await userStore.reloadUserAfterRefresh();
    if (currentUser.value && currentUser.value.id) {
      await followStore.fetchFollowing(currentUser.value.id);
      await groupStore.fetchUserGroups();
      removeWsListener.value = wsService.addListener('chat-component', handleIncomingMessage);
    }
  } catch (error) {
    console.error("Error during component mount:", error);
  }
});

onUnmounted(() => {
  if (removeWsListener.value) {
    removeWsListener.value();
  }
});

</script>

<style scoped>
.chat-container { 
  display: flex; 
  height: 100vh; 
  background-color: #0a0a0f; 
  color: #d1d5db;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
}
.chat-layout { 
  display: flex; 
  flex-direction: row; 
  height: 100%; 
  width: 100%; 
  overflow: hidden; 
}

.friends-list-section { 
  display: flex; 
  flex-direction: column; 
  padding: 1.5rem; 
  width: 24rem; 
  background: #181824; 
  border-right: 1px solid rgba(255, 255, 255, 0.1); 
  flex-shrink: 0; 
}
.header { 
  display: flex; 
  align-items: center; 
  gap: 0.75rem;
  height: 3rem; 
  width: 100%; 
  flex-shrink: 0; 
  margin-bottom: 1.5rem; 
}
.header-icon-wrapper { 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  border-radius: 0.75rem; 
  color: #000; 
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%); 
  height: 2.5rem; 
  width: 2.5rem; 
}
.header-icon { width: 1.5rem; height: 1.5rem; }
.header-title { 
  font-weight: 700; 
  font-size: 1.5rem; 
  background: linear-gradient(135deg, #ffffff 0%, #e879c6 50%, #78c7ff 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.profile-card { 
  display: flex; 
  flex-direction: column; 
  align-items: center; 
  background: rgba(15, 15, 23, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1); 
  margin-top: 1rem; 
  width: 100%; 
  padding: 1.5rem 1rem; 
  border-radius: 1rem; 
  flex-shrink: 0;
  margin-bottom: 1.5rem;
}
.profile-avatar-wrapper { 
  height: 5rem; 
  width: 5rem; 
  border-radius: 50%; 
  border: 3px solid #e879c6;
  padding: 3px;
  background: #0a0a0f;
  overflow: hidden;
}
.profile-avatar { height: 100%; width: 100%; object-fit: cover; border-radius: 50%; }
.profile-name { font-size: 1.1rem; font-weight: 600; margin-top: 0.75rem; color: #ffffff; }
.profile-status { font-size: 0.875rem; color: #10b981; }

.friends-list-container { display: flex; flex-direction: column; margin-top: 1.5rem; flex-grow: 1; min-height: 0; }
.friends-list-header { display: flex; justify-content: space-between; align-items: center; font-size: 0.875rem; padding: 0 0.5rem; color: #9ca3af; }
.friends-list-title { font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; }
.friends-list-count { 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  background-color: #374151; 
  height: 1.25rem; 
  min-width: 1.25rem; 
  padding: 0 0.5rem; 
  border-radius: 0.625rem; 
  font-size: 0.75rem; 
  font-weight: bold; 
  color: #d1d5db; 
}
.friends-list { display: flex; flex-direction: column; gap: 0.5rem; margin-top: 1rem; flex-grow: 1; overflow-y: auto; padding-right: 0.5rem; }
.friend-item { 
  display: flex; 
  align-items: center; 
  border-radius: 0.75rem; 
  padding: 0.75rem; 
  width: 100%; 
  text-align: left; 
  border: 1px solid transparent; 
  background: none; 
  cursor: pointer; 
  transition: all 0.3s ease; 
  position: relative; 
  color: #d1d5db; 
}
.friend-item:hover { background-color: rgba(255, 255, 255, 0.05); }
.friend-item.active { 
  color: #ffffff; 
  background: linear-gradient(135deg, rgba(232, 121, 198, 0.2) 0%, rgba(120, 199, 255, 0.2) 100%); 
  font-weight: 600; 
  border: 1px solid rgba(232, 121, 198, 0.3);
}
.friend-avatar-wrapper { height: 3rem; width: 3rem; border-radius: 50%; flex-shrink: 0; overflow: hidden; margin-right: 0.75rem; border: 2px solid #374151; }
.friend-item.active .friend-avatar-wrapper { border-color: #e879c6; }
.friend-avatar { height: 100%; width: 100%; object-fit: cover; }
.friend-details { display: flex; flex-direction: column; }
.friend-name { font-size: 1rem; font-weight: 500; color: inherit; }
.friend-item.active .friend-name { font-weight: 600; color: #ffffff; }
.friend-last-message { font-size: 0.8rem; color: #9ca3af; }
.unread-badge {
  position: absolute;
  top: 5px;
  right: 5px;
  background: #e879c6;
  color: #000;
  border-radius: 50%;
  height: 1.5rem;
  width: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  font-weight: bold;
  border: 2px solid #181824;
}

.chat-area-wrapper { display: flex; flex-direction: column; flex: 1 1 0%; height: 100%; padding: 1.5rem; }
.chat-box { display: flex; flex-direction: column; flex: 1 1 0%; border-radius: 1.5rem; background: #181824; border: 1px solid rgba(255, 255, 255, 0.1); height: 100%; min-height: 0; }
.active-chat { display: flex; flex-direction: column; height: 100%; }
.chat-header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; border-bottom: 1px solid rgba(255, 255, 255, 0.1); flex-shrink: 0; }
.chat-header-user { display: flex; align-items: center; gap: 1rem; }
.chat-header-avatar { width: 3rem; height: 3rem; border-radius: 50%; object-fit: cover; border: 2px solid #e879c6;}
.chat-header-text { display: flex; flex-direction: column; line-height: 1.25; }
.chat-header-name { font-size: 1.25rem; font-weight: 600; color: #ffffff; }
.chat-header-status { font-size: 0.875rem; color: #10b981; }

.messages-container { display: flex; flex-direction: column; flex-grow: 1; overflow-y: auto; padding: 1.5rem; min-height: 0; }
.messages-list { margin-top: auto; }
.messages-grid { display: flex; flex-direction: column; gap: 1rem; }
.message-wrapper { display: flex; max-width: 75%; }
.message-wrapper.received { align-self: flex-start; }
.message-wrapper.sent { align-self: flex-end; }
.message-layout { display: flex; align-items: flex-end; flex-direction: row; gap: 0.75rem; }
.message-layout.sent { flex-direction: row-reverse; }
.message-avatar-wrapper { height: 2.5rem; width: 2.5rem; border-radius: 50%; flex-shrink: 0; overflow: hidden; }
.message-avatar { width: 100%; height: 100%; object-fit: cover; }
.message-bubble { position: relative; font-size: 1rem; padding: 0.75rem 1rem; border-radius: 1.25rem; color: #d1d5db; word-break: break-word; }
.message-bubble.received { background-color: #23233a; border-top-left-radius: 0.25rem; }
.message-bubble.sent { background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%); color: #000; font-weight: 500; border-top-right-radius: 0.25rem; }

.message-input-section { display: flex; align-items: center; padding: 1rem 1.5rem; border-top: 1px solid rgba(255, 255, 255, 0.1); flex-shrink: 0; gap: 1rem; position: relative; background: #181824; border-bottom-left-radius: 1.5rem; border-bottom-right-radius: 1.5rem;}
.emojis-container { position: relative; }
.emoji-toggle-btn { background: none; border: none; font-size: 1.5rem; cursor: pointer; color: #9ca3af; padding: 0.5rem; transition: all 0.3s ease; }
.emoji-toggle-btn:hover { color: #e879c6; transform: scale(1.1); }
.emojis-picker { position: absolute; bottom: 100%; left: 0; background: #23233a; border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 1rem; padding: 0.5rem; display: flex; flex-wrap: wrap; gap: 0.5rem; width: 200px; box-shadow: 0 -4px 20px rgba(0,0,0,0.2); z-index: 10; margin-bottom: 0.5rem; }
.emoji-item { cursor: pointer; font-size: 1.25rem; transition: transform 0.2s; }
.emoji-item:hover { transform: scale(1.2); }
.message-input { 
  flex-grow: 1; 
  background-color: rgba(255, 255, 255, 0.05); 
  border: 1px solid rgba(255, 255, 255, 0.1); 
  border-radius: 1rem; 
  padding: 0 1rem; 
  height: 3.5rem; 
  color: #ffffff; 
  font-size: 1rem; 
  transition: all 0.3s; 
}
.message-input::placeholder { color: rgba(255, 255, 255, 0.5); }
.message-input:focus { 
  outline: none; 
  border-color: #e879c6; 
  background: rgba(255, 255, 255, 0.1);
  box-shadow: 0 0 0 3px rgba(232, 121, 198, 0.2);
}
.send-button { 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%); 
  border-radius: 1rem; 
  color: #000; 
  height: 3.5rem; 
  width: 3.5rem; 
  flex-shrink: 0; 
  border: none; 
  cursor: pointer; 
  transition: all 0.3s ease; 
}
.send-button:hover { 
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(232, 121, 198, 0.4);
}
.send-button-icon { width: 1.5rem; height: 1.5rem; }

.empty-chat-container { display: flex; flex-direction: column; height: 100%; align-items: center; justify-content: center; text-align: center; padding: 2rem; }
.empty-chat-icon-wrapper { 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  border-radius: 50%; 
  color: #000; 
  background: linear-gradient(135deg, #e879c6 0%, #78c7ff 100%); 
  height: 5rem; 
  width: 5rem; 
}
.empty-chat-icon { width: 2.5rem; height: 2.5rem; }
.empty-chat-title { font-size: 1.5rem; font-weight: 600; margin-top: 1.5rem; color: #ffffff; }
.empty-chat-text { color: #9ca3af; margin-top: 0.5rem; }

.custom-scrollbar::-webkit-scrollbar { width: 8px; }
.custom-scrollbar::-webkit-scrollbar-track { background: #181824; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #374151; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #4b5563; }

@media (max-width: 900px) {
  .chat-layout {
    flex-direction: column;
  }
  .friends-list-section {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    height: auto;
  }
  .chat-area-wrapper {
    padding: 1rem;
  }
}

@media (max-width: 600px) {
  .friends-list-section {
    padding: 1rem;
  }
  .header {
    margin-bottom: 1rem;
  }
  .profile-card {
    display: none;
  }
  .chat-area-wrapper {
    padding: 0.5rem;
  }
  .chat-box {
    border-radius: 0;
  }
  .chat-header {
    padding: 0.75rem;
  }
  .message-input-section {
    padding: 0.75rem;
    gap: 0.5rem;
  }
  .message-input, .send-button {
    height: 3rem;
  }
  .send-button {
    width: 3rem;
  }
}
.chat-container { 
  display: flex; 
  /* Changed height to be dynamic */
  height: calc(100vh - 60px); /* Full height minus top navbar */
  background-color: #0a0a0f; 
  color: #d1d5db;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', sans-serif;
}
.chat-layout { 
  display: flex; 
  flex-direction: row; 
  height: 100%; 
  width: 100%; 
  overflow: hidden; 
}

/* ... other base styles ... */
.friends-list-section { 
  display: flex; 
  flex-direction: column; 
  padding: 1.5rem; 
  width: 24rem; 
  background: #181824; 
  border-right: 1px solid rgba(255, 255, 255, 0.1); 
  flex-shrink: 0; 
  height: 100%;
  overflow-y: auto;
}
.chat-area-wrapper { 
  display: flex; 
  flex-direction: column; 
  flex: 1 1 0%; 
  height: 100%; 
  padding: 0; /* Remove padding to make chat-box fill the area */
}
.chat-box { 
  display: flex; 
  flex-direction: column; 
  flex: 1 1 0%; 
  border-radius: 0; /* Remove radius for seamless view */
  background: #181824; 
  border: none;
  height: 100%; 
  min-height: 0; 
}
.messages-container { 
  display: flex; 
  flex-direction: column; 
  flex-grow: 1; 
  overflow-y: auto; 
  padding: 1.5rem; 
  min-height: 0; 
}


/* --- Responsive Fixes --- */
@media (max-width: 900px) {
  .chat-container {
    /* Full viewport height minus top and bottom navs on mobile */
    height: calc(100vh - 120px);
  }
  .chat-layout {
    flex-direction: column;
  }
  .friends-list-section {
    width: 100%;
    height: auto; /* Auto height for the top bar */
    flex-shrink: 0;
    border-right: none;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    overflow-y: hidden; /* Prevent vertical scroll */
  }
  .friends-list-container {
     flex-shrink: 0; /* Prevent list from shrinking */
  }
  .friends-list {
    flex-direction: row; /* Horizontal scroll for lists */
    overflow-x: auto;
    padding-bottom: 1rem; /* Space for scrollbar */
  }
  .friend-item {
    flex-direction: column;
    min-width: 80px; /* Give items more space */
    gap: 8px;
  }
  .friend-avatar-wrapper {
    width: 50px;
    height: 50px;
    margin-right: 0;
  }
  .friend-name {
    font-size: 12px;
    white-space: nowrap;
  }
  .profile-card {
    display: none; /* Hide profile card on mobile to save space */
  }
  .header-title {
      font-size: 1.2rem;
  }
  .chat-area-wrapper {
      padding: 0;
      flex-grow: 1;
      min-height: 0; /* Allow chat area to shrink and grow */
  }
}

@media (max-width: 600px) {
  .chat-header {
    padding: 0.75rem;
  }
  .messages-container {
    padding: 0.75rem;
  }
  .message-input-section {
    padding: 0.75rem;
    gap: 0.5rem;
  }
}
</style>