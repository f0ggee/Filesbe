## This file contains the explanation about some secure elements of my app



<img width="4308" height="3252" alt="image" src="https://github.com/user-attachments/assets/cf32da42-09f9-4298-9db8-176d3be160bc" />



1)The scheme starts with generating a key which it will send to other servers(the main notes, we don't generate keys, only a key, this is a  mistake)


2)Then we generate a hash and sign it with a  master server's key. This has been added so that a server which will get a packet make sure that it is from the Master server truly 


<img width="6324" height="4553" alt="image" src="https://github.com/user-attachments/assets/23874f37-b1a8-4b2b-b281-caa699d1eb5e" />

