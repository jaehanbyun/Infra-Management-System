package app

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/bndr/gojenkins"
	"github.com/gophercloud/gophercloud"
)

func getAuthOpts() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: os.Getenv("OS_AUTH_URL"),
		Username:         os.Getenv("OS_USERNAME"),
		Password:         os.Getenv("OS_PASSWORD"),
		DomainName:       os.Getenv("OS_USER_DOMAIN_NAME"),
	}
}

func triggerJenkinsJob(jobName string, params map[string]string) error {
	ctx := context.Background()
	jenkins, err := gojenkins.CreateJenkins(nil, JenkinsURL, JenkinsUser, JenkinsAPIToken).Init(ctx)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	buildId, err := jenkins.BuildJob(ctx, jobName, params)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Job %s started with build number %d\n", jobName, buildId))
	}

	return nil
}

func structToMap(item interface{}) map[string]string {
	result := make(map[string]string)
	rValue := reflect.ValueOf(item)
	rType := reflect.TypeOf(item)

	for i := 0; i < rType.NumField(); i++ {
		fieldValue := rValue.Field(i)
		fieldType := rType.Field(i)

		key := fieldType.Tag.Get("json")
		if key == "" {
			key = fieldType.Name
		}

		switch fieldValue.Kind() {
		case reflect.String:
			result[key] = fieldValue.String()
		case reflect.Int:
			result[key] = fmt.Sprintf("%d", fieldValue.Int())
		}
	}

	return result
}
