-- sqlite3 -init imageboard.sql
	 
			PRAGMA journal_mode = WAL;
			PRAGMA synchronous = NORMAL;
			PRAGMA temp_store = MEMORY;
			 
			BEGIN;
			 
			CREATE TABLE IF NOT EXISTS post (
			        id INTEGER PRIMARY KEY AUTOINCREMENT,
			        subject TEXT DEFAULT null,
			        name TEXT,
			        text TEXT,
			        date_posted TEXT DEFAULT CURRENT_TIMESTAMP,
			        thread_id INTEGER
			);
			 
			CREATE INDEX IF NOT EXISTS idx_thread_id ON post (thread_id, date_posted);
			 
			CREATE VIEW IF NOT EXISTS latest_threads AS
			        SELECT post.thread_id, post.subject, post.name, post.text, post.date_posted
			        FROM post
			        ORDER BY post.thread_id DESC, post.date_posted DESC;
			 
			-- only keep 5 threads
			CREATE TRIGGER purge_old_posts
			AFTER INSERT ON post
			        BEGIN
			                DELETE FROM post
			                WHERE (
			                        SELECT thread_id  
			                        FROM latest_threads
			                        GROUP BY thread_id
			                        ORDER BY thread_id DESC  
			                        LIMIT 10 OFFSET 5
			                ) = post.thread_id;
			        END;

			COMMIT;
