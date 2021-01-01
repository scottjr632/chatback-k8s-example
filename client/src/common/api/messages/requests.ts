import axios from 'axios';
import { Message } from '../../types';

export async function getAllMessages(): Promise<Message[]> {
  const { data } = await axios.get('/api/messages');
  return data;
}