
#### Sending out emails from your api server 
------
More often you would need your APIs sending some email, specifically notifications to intended users. This requires the api to have access to sender's email using SMTP. 

#### Setting up your email account for this
----
Before you get along with sending emails from api handlers, you need to ensure some settings on your mailbox. 

##### Gmail settings 
------
In the context of gmail lets see what settings are required. Since May2022 `Less secured apps` access is completely removed and 2FA is mandated. With 2FA you need to generate an app password that you can specifically use only when you access emails from your APIs. App password is just text string that can be vulnerable when used directly inside the code. Remember to use environment files and load the passwords from within the code. (Instead of hardcoding the passwords as naked text)

#### Text/HTML mails 
------

#### Emals with attachments 
-------

#### Timeouts and using co-routines
-------
