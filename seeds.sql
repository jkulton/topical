INSERT INTO topics (title)
VALUES ('Super simple file server?'),
       ('Found a Spotify playlist that is :fire:'),
       ('Share your best NBA gifs!');

INSERT INTO messages (topic_id, author_initials, author_theme, content)
VALUES (1, 'BT', 4, 'What is a simple tool/command for starting a file server fast?

Say I''m just trying to serve a few files in my local directory over HTTP, or I wanna serve a simple HTML file?

How would I do that easily'),
       (1, 'JK', 1, 'Check out [Platter](https://github.com/jkulton/platter).

It''s a simple program written in Go that let''s you serve your current directory (or one you specify) over HTTP, fast.

If you have `$GOPATH/bin` in your `$PATH` just run:

```
go get github.com/jkulton/platter/cmd/platter
```

and then

```
platter
```

to serve your current working directory over HTTP.'),
       (2, 'YM', 3, 'Check out this [Instrumental Lofi playlist on Spotify](https://open.spotify.com/user/spotify/playlist/37i9dQZF1DXc8kgYqQLMfH?si=IfYoL7gyTp-75RCbeEc_9A).
I''ve been digging so many songs off this.'),
       (2, 'JK', 1, 'I LOVE THIS. Thanks for sharing!'),
       (2, 'HK', 4, 'Hey I know this playlist! It''s one of my favs.'),
       (2, 'YM', 3, '@HK, good taste! :)'),
       (2, 'KB', 5, '![10/10](https://i.imgur.com/06bFgD2.gif)

Love it, thx for sharing!'),
       (3, 'DL', 2, 'Share your best NBA gifs!
Bring them on!

![Alt text](https://i.imgur.com/ao5tOjM.gif)'),
       (3, 'MB', 6, '![Kobe doesn''t flinch](https://i.imgur.com/uOPgcGJ.gif)'),
       (3, 'JH', 2, '@MB nice try bro **LOL**

![James Harden](https://i.imgur.com/X0yWmKZ.gif)'),
       (3, 'KI', 4, '![Kyrie''s big shot](https://i.imgur.com/IM8y2mL.gif)'),
       (3, 'KB', 5, 'tfw someone says hard tacos are better

![Alt text](https://i.imgur.com/HPnZFlO.gif)');