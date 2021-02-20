package mysql

import (
	"database/sql"

	"arnoldcodes.com/snippetbox/pkg/models"
)

//Define a SnippetModel type which wraps a sql.FB connection pool
type SnippetModel struct {
	DB *sql.DB
}

//This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	//sql statement we want to execute
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	//Use the Exec() method on the embedded connection pool to execute the statment. The first parameter is the SQL statement, followed by the title, content and expiry values for the placeholder parameters.
	//This method returns a sql.Result object, which contains some basic information about what happened when the statement was executed
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	//Use the LastInsertId() methos on the result object to get the ID of our
	//newly inserted record in the snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	//The ID returned has the type int64, so we convert it to an itn type
	//before returning.
	return int(id), nil

}

//Getting a specific snippet by its id
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	row := m.DB.QueryRow(stmt, id)
	//Initialize a pointer to a new zeroed Snippet struct
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

//This will return The 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	//SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`
	//User the Query() method on the connection pool to execute our
	//SQL statement.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	//We defer rows.Close() to ensure the sql.Rows resultset is
	//always properly closed before the Latest() method returns. This defer statement should
	//come *after* you check for an error from the Query() method.
	//Otherwise, if query() returns an error, you'll get a panic
	//trying to close a nil resultset
	defer rows.Close()

	//Inirialize an empty slice to hold the models.Snippets objects.
	snippets := []*models.Snippet{}

	//Use rows.Next to iterate through the rows in the resultset. This
	//prepares the first (and then each subsequent) row to be acted on by the rows.Scan() method.
	//if iteration over all the rows completes then the resultset automatically closes itself and frees-up
	//the underlying database connection
	for rows.Next() {
		//create a pointer to a new zeroed Snippet struct
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}
		//append it to the slice of snippets
		snippets = append(snippets, s)

	}
	//When the rows.Next() loop has finished we call rows.Err() to retrieve an
	//error that was encountered duringthe iteration. It's important to call this- don't assume
	//that a succefull iteration was completed
	//over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
