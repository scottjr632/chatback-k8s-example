import React, { FC, useContext, useEffect, useRef, useState } from 'react';

import { Message } from '../types';
import ReconnectingWebSocket from 'reconnecting-websocket';

interface State {
  message: Message | undefined;
  sendMessage: (send: Message) => any;
}

const MessagesContext = React.createContext<State>({
  message: undefined,
  sendMessage: () => null,
});

export const MessagesProvider: FC = ({ children }) => {
  const socketRef = useRef<ReconnectingWebSocket | null>(null);
  const [message, setMessage] = useState<Message | undefined>();

  useEffect(() => {
    socketRef.current = new ReconnectingWebSocket(`ws://${window.location.host}/ws/messages`);
    socketRef.current.addEventListener('message', e => {
      console.log(e);
      setMessage(JSON.parse(e.data));
    });
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