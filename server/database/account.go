package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"templateapp/models"
)

func FindAccountByHandleNameAndPassword(handleName string, passwordString string) (*models.Account, error) {
	var account models.Account

	allAccount := fmt.Sprintf("%saccount.*", TablePrefix)
	joinPassword := fmt.Sprintf("JOIN %slogin_password_grant lpg ON lpg.login_handle_name_grant_id = lhng.id", TablePrefix)
	joinHandleName := fmt.Sprintf("JOIN %slogin_handle_name_grant lhng ON lhng.account_id = %saccount.id", TablePrefix, TablePrefix)

	if find := DB().Model(&account).Select(allAccount).Joins(joinHandleName).Joins(joinPassword).Where("lhng.handle_name = ?", handleName).Where("lpg.password_string = ?", passwordString).First(&account); nil != find.Error {
		if errors.Is(find.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
	}
	return &account, nil
}
