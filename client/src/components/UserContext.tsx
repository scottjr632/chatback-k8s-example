import { Modal, ModalOverlay, ModalContent, ModalHeader, ModalCloseButton, ModalBody, FormControl, FormLabel, Input, ModalFooter, Button, useDisclosure } from '@chakra-ui/react';
import React, { FC, useCallback, useContext, useEffect, useState } from 'react';

interface Props {
  storageKey: string;
}

interface State {
  username: string;
  hasInitialized: boolean;
  clear: () => any;
  openModal: () => any;
  setUsername: (_: string) => any;
}

export const UserContext = React.createContext<State>({
  username: '',
  hasInitialized: false,
  clear: () => null,
  openModal: () => null,
  setUsername: () => null,
});
const UserProvider: FC<Props> = ({ children, storageKey }) => {
  const [username, setUsername] = useState('');
  const [hasInitialized, setHasInitialized] = useState(false);

  const { isOpen, onOpen, onClose } = useDisclosure({defaultIsOpen: false});

  useEffect(() => {
    const cachedUsername = localStorage.getItem(storageKey);
    if (cachedUsername)
      setUsername(cachedUsername);
    else
      onOpen();

    setHasInitialized(true);
  }, [onOpen, storageKey]);

  const setUsernameInCache = useCallback((username: string) => {
    localStorage.setItem(storageKey, username);
    setUsername(username);
  }, [storageKey]);

  const clear = useCallback(() => {
    localStorage.removeItem(storageKey);
    setUsername('');
  }, [storageKey]);

  return (
    <UserContext.Provider value={{
      username,
      hasInitialized,
      clear,
      openModal: onOpen,
      setUsername: setUsernameInCache,
    }}>
      <Modal
        isOpen={isOpen}
        onClose={onClose}
      >
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>What's your name?</ModalHeader>
          <ModalCloseButton />
          <ModalBody pb={6}>
            <FormControl>
              <FormLabel>Username</FormLabel>
              <Input placeholder="Username" value={username} onChange={e => {
                setUsername(e.currentTarget.value);
              }} />
            </FormControl>
          </ModalBody>

          <ModalFooter>
            <Button colorScheme="blue" mr={3} onClick={() => {
              setUsernameInCache(username);
              onClose();
            }}>
              Save
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => useContext(UserContext);
export default UserProvider;