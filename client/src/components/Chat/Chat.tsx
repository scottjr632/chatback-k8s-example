import React, { useEffect, useRef, useState } from 'react';
import { Box, Button, Flex, Input } from '@chakra-ui/react';

import { useUser } from '../UserContext';
import { Message as MesageIntf } from '../../common/types';
import { useMessages } from '../../common/api/messages/hooks';
import { useMessageSocket } from '../../common/api/MessagesProvider';

import Message from './Message';

interface Props {
  messages?: MesageIntf[];
}

export default function Chat() {
  const { username } = useUser();
  const { data } = useMessages();
  const { message, sendMessage } = useMessageSocket();

  const [content, setContent] = useState('');
  const [messages, setMessages] = useState<MesageIntf[]>([]);

  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (data)
      setMessages(data);
  }, [data]);

  useEffect(() => {
    if (message) {
      setMessages(prevMessages => [...prevMessages, message]);
    }
  }, [message]);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  const send = () => {
    if (content !== '') {
      sendMessage({
        username,
        content,
        created: new Date().toString(),
      });
      setContent('');
    }
  };

  const handleEnter = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter')
      send();
  };

  return (
    <Flex
      maxW='1920px'
      maxH='100%'
      minH='100%'
      marginX='auto'
      direction='column'
      overflow='scroll'
    >
      <Box height='100%' overflow='scroll' flex='1 1 auto' ref={scrollRef}>
        {messages.map(message => (
          <Box
            width='90%'
            marginX='auto'
            marginY='1rem'
          >
            <Message {...message} />
          </Box>
        ))}
      </Box>
      {username &&
      <Flex padding='0.5rem'>
        <Input
          value={content}
          onKeyPress={handleEnter}
          placeholder='Press enter to send'
          onChange={e => setContent(e.currentTarget.value)}
        />
        <Button onClick={() => {
          sendMessage({
            username,
            content,
            created: new Date().toString(),
          });
          setContent('');
        }}>Send</Button>
      </Flex>}
    </Flex>
  );
}