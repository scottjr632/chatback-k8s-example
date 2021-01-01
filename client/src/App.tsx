import React from 'react';
import {
  QueryClient,
  QueryClientProvider,
} from 'react-query';
import {
  Box,
  ChakraProvider,
  Flex,
  theme,
} from '@chakra-ui/react';

import Chat from './components/Chat';
import Header from './components/Header';
import UserProvider from './components/UserContext';
import MessagesProvider from './common/api/MessagesProvider';

const queryClient = new QueryClient();

export const App = () => (
  <ChakraProvider theme={theme}>
    <QueryClientProvider client={queryClient}>
      <UserProvider storageKey='chatback-username'>
        <MessagesProvider>
          <Flex flexDirection='column' height='100vh'>
            <Header title='ChatBack' />
            <Box height='100%' flex='1 1 auto' overflow='scroll'>
              <Chat />
            </Box>
          </Flex>
        </MessagesProvider>
      </UserProvider>
    </QueryClientProvider>
  </ChakraProvider>
);
