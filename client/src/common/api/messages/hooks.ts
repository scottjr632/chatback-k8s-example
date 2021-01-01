import { useQuery } from 'react-query';
import { getAllMessages } from './requests';

export function useMessages() {
  return useQuery('messages', getAllMessages);
}