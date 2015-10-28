# Setting up a Chrome runner

**Run in remote debugging mode, using a custom user data directory:**

    chrome --remote-debugging-port=9222 --user-data-dir=./data

This allows us to set content preferences -- in particular, running with images disabled.
After the first run, preferences live in `[user-data-dir]/Default/Preferences`. 

As of Chrome v.46, this is a json file so it is easy to inspect and manipulate. The setting that we are looking for is `profile.default_content_setting_values.images = 2`