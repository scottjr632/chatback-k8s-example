import React, { useMemo } from 'react';
import { Avatar, Box, Flex, Text, useColorModeValue } from '@chakra-ui/react';

import { Message as MessageIntf } from '../../common/types';

type Props = MessageIntf;

export default function Message({
  username,
  avatarUrl,
  content,
  created: date,
}: Props) {
  const bg = useColorModeValue('rgba(0,0,0,0.1)', 'rgba(0,0,0,0.4)');
  const localDate = useMemo(() => {
    const parsedDate = new Date(date);
    return `${parsedDate.toLocaleDateString()} ${parsedDate.toLocaleTimeString()}`;
  }, [date]);

  return (
    <Box bg={bg} borderRadius='8px' padding='1rem 0.5rem'>
      <Flex paddingX={'1.2rem'}>
        <Avatar name={username} src={avatarUrl} />
        <Box marginY='auto' marginX='1rem'>
          <Text fontWeight='bold'>{username}</Text>
          <Text>
            {localDate}
          </Text>
        </Box>
      </Flex>
      <Text
        paddingX={'2rem'} 
        paddingY={'0.5rem'}
      >
        {content}
      </Text>
    </Box>
  );
}