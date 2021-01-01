import React from 'react';
import { Box, Flex, Grid, Heading, Link, Text } from '@chakra-ui/react';

import { ColorModeSwitcher } from './ColorModeSwitcher';
import { useUser } from './UserContext';

interface Props {
  title: string;
  username?: string;
}

export default function Header({ title }: Props) {
  const { username, openModal, clear } = useUser();
  return (
    <Grid templateColumns='1fr 1fr' width='100%' paddingX='1rem' paddingY='0.5rem'>
      <Box width='100%'>
        <Heading>{title}</Heading>
      </Box>
      <Flex width='100%' justifyContent='flex-end'>
        <Flex marginY='auto'>
          <Text fontSize='lg'>
            {username}
            <Link onClick={username ? clear : openModal} marginLeft='1rem'>
              {username ? 'Log out' : 'Set username'}
            </Link>
          </Text>
        </Flex>
        <ColorModeSwitcher justifySelf="flex-end" />
      </Flex>
    </Grid>
  );
}