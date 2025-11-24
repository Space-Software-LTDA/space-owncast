import { createContext } from 'react';
import { ChatMessage } from '../interfaces/chat-message.model';
import { getUnauthedData } from '../utils/apis';

const ENDPOINT = `/api/chat`;
const URL_CHAT_REGISTRATION = `/api/chat/register`;

export interface UserRegistrationResponse {
  id: string;
  accessToken: string;
  displayName: string;
  displayColor: number;
}

export interface ChatStaticService {
  getChatHistory(accessToken: string): Promise<ChatMessage[]>;
  registerUser(username: string, email: string): Promise<UserRegistrationResponse>;
}

class ChatService {
  public static async getChatHistory(accessToken: string): Promise<ChatMessage[]> {
    try {
      const response = await getUnauthedData(`${ENDPOINT}?accessToken=${accessToken}`);
      return response;
    } catch (e) {
      console.error(e);
      return [];
    }
  }

  public static async registerUser(
    username: string,
    email: string,
  ): Promise<UserRegistrationResponse> {
    const payload: { displayName: string; email?: string } = { displayName: username };

    // Only include email if it's not empty
    if (email && email.trim() !== '') {
      payload.email = email;
    }

    console.log('Registering user with payload:', payload);
    const options = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      data: payload,
    };

    const response = await getUnauthedData(URL_CHAT_REGISTRATION, options);

    console.log('User registration response:', response);
    return response;
  }
}

export const ChatServiceContext = createContext<ChatStaticService>(ChatService);
