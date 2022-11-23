package repositories

import (
	"github.com/sayopaul/sendchamp-go-test/infrastructure"
	"github.com/sayopaul/sendchamp-go-test/models"
)

type TaskRepository struct {
	db infrastructure.Database
}

func NewTaskRepository(db infrastructure.Database) TaskRepository {
	return TaskRepository{
		db: db,
	}
}

func (tr TaskRepository) FetchTask(taskCondition models.Task) (task models.Task, err error) {
	return task, tr.db.DB.Model(&models.Task{}).Where(&taskCondition).First(&task).Error
}

func (tr TaskRepository) CreateTask(task models.Task) (models.Task, error) {
	return task, tr.db.DB.Create(&task).Error
}

func (tr TaskRepository) UpdateTask(taskCondition models.Task, task models.Task) (updatedTask models.Task, err error) {
	err = tr.db.DB.Model(models.Task{}).Where(&taskCondition).Updates(&task).Error
	tr.db.DB.Model(&models.Task{}).Where(&task).First(&updatedTask)
	return updatedTask, err

}

func (tr TaskRepository) DeleteTask(taskCondition models.Task) error {
	return tr.db.DB.Where(&taskCondition).Delete(&models.Task{}).Error
}
