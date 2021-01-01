import React, { FC, useContext, useEffect, useRef, useState } from 'react';
import { Message } from '../types';

interface State {
  message: Message | undefined;
  sendMessage: (send: Message) => any;
}

const MessagesContext = React.createContext<State>({
  message: undefined,
  sendMessage: () => null,
});

export const MessagesProvider: FC = ({ children }) => {
  const socketRef = useRef<WebSocket | null>(null);
  const [message, setMessage] = useState<Message | undefined>();

  useEffect(() => {
    socketRef.current = new WebSocket(`ws://${window.location.host}/ws/messages`);
    socketRef.current.onmessage = (e) => {
      console.log(e);
      setMessage(JSON.parse(e.data));
    };
  }, []);

  const sendMessage = (send: Message) => {
    if (socketRef.current)
      socketRef.current.send(JSON.stringify(send));
  };

  return (
    <MessagesContext.Provider value={{
      message,
      sendMessage,
    }}>
      {children}
    </MessagesContext.Provider>
  );
};

export const useMessageSocket = () => useContext(MessagesContext);
export default MessagesProvider;