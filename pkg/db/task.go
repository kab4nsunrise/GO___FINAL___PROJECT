package db

import (
    "database/sql"
    "fmt"
)

type Task struct {
    ID      string `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
    res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
        task.Date, task.Title, task.Comment, task.Repeat)
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
    rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    tasks := []*Task{}
    for rows.Next() {
        t := &Task{}
        err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, t)
    }
    return tasks, nil
}

func GetTask(id string) (*Task, error) {
    row := db.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)
    t := &Task{}
    err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("задача не найдена")
    }
    return t, err
}

func UpdateTask(task *Task) error {
    res, err := db.Exec(`UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`,
        task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return err
    }
    cnt, _ := res.RowsAffected()
    if cnt == 0 {
        return fmt.Errorf("задача не найдена")
    }
    return nil
}

func DeleteTask(id string) error {
    res, err := db.Exec(`DELETE FROM scheduler WHERE id=?`, id)
    if err != nil {
        return err
    }
    cnt, _ := res.RowsAffected()
    if cnt == 0 {
        return fmt.Errorf("задача не найдена")
    }
    return nil
}

func UpdateDate(id, newDate string) error {
    res, err := db.Exec(`UPDATE scheduler SET date=? WHERE id=?`, newDate, id)
    if err != nil {
        return err
    }
    cnt, _ := res.RowsAffected()
    if cnt == 0 {
        return fmt.Errorf("задача не найдена")
    }
    return nil
}

func SearchTasks(search string, limit int) ([]*Task, error) {
    pattern := "%" + search + "%"
    rows, err := db.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
        pattern, pattern, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    tasks := []*Task{}
    for rows.Next() {
        t := &Task{}
        rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
        tasks = append(tasks, t)
    }
    return tasks, nil
}
