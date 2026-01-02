package tests

import (
	"context"

	"github.com/cucumber/godog"
)

// DepartmentTestContext holds the test context for department scenarios
type DepartmentTestContext struct {
	data map[string]interface{}
}

// NewDepartmentTestContext creates a new test context
func NewDepartmentTestContext() *DepartmentTestContext {
	return &DepartmentTestContext{
		data: make(map[string]interface{}),
	}
}

// InitializeScenario sets up godog step definitions
func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := NewDepartmentTestContext()

	// Background steps
	ctx.Step(`^the Department Management API is running$`, tc.apiIsRunning)
	ctx.Step(`^the Redis cache is available$`, tc.redisIsAvailable)
	ctx.Step(`^the database is initialized$`, tc.databaseIsInitialized)

	// Given steps - Setup
	ctx.Step(`^I have department data with code "([^"]*)" and name "([^"]*)"$`, tc.haveDepartmentData)
	ctx.Step(`^the annual budget is \$([0-9,]+)$`, tc.setAnnualBudget)
	ctx.Step(`^the department head is "([^"]*)"$`, tc.setDepartmentHead)
	ctx.Step(`^the status is "([^"]*)"$`, tc.setStatus)
	ctx.Step(`^no status is specified$`, tc.noStatusSpecified)

	// Given steps - Existing data
	ctx.Step(`^a department with code "([^"]*)" exists in the system$`, tc.departmentExists)
	ctx.Step(`^a department with ID (\d+) exists in the system$`, tc.departmentWithIDExists)
	ctx.Step(`^a department with code "([^"]*)" already exists$`, tc.departmentAlreadyExists)
	ctx.Step(`^a department with code "([^"]*)" exists$`, tc.departmentExists)
	ctx.Step(`^a department with code "([^"]*)" and status "([^"]*)" exists$`, tc.departmentWithStatusExists)

	// Given steps - Cache state
	ctx.Step(`^the Redis cache is empty for "([^"]*)"$`, tc.cacheEmpty)
	ctx.Step(`^a department with code "([^"]*)" is cached in Redis$`, tc.cachedInRedis)
	ctx.Step(`^"([^"]*)" is not in Redis cache$`, tc.notInCache)
	ctx.Step(`^a department with code "([^"]*)" was cached (\d+) minutes ago$`, tc.cachedMinutesAgo)

	// Given steps - Multiple departments
	ctx.Step(`^the following departments exist:$`, tc.multipleDepartmentsExist)
	ctx.Step(`^multiple departments exist with different statuses$`, tc.differentStatuses)
	ctx.Step(`^multiple departments exist with different budgets and utilization$`, tc.differentBudgets)

	// Given steps - Budget
	ctx.Step(`^I create a department with code "([^"]*)" and budget \$([0-9,]+)$`, tc.createDepartmentWithBudget)
	ctx.Step(`^a department with code "([^"]*)" has annual budget \$([0-9,]+)$`, tc.departmentHasBudget)
	ctx.Step(`^the department has spent \$([0-9,]+)$`, tc.departmentSpent)

	// Given steps - Invalid data
	ctx.Step(`^I have department data with code "([^"]*)" \(lowercase\)$`, tc.lowercaseCode)
	ctx.Step(`^I have department data with code "([^"]*)" \((\d+) letter[s]?\)$`, tc.codeWithLetters)
	ctx.Step(`^I have department data with valid code "([^"]*)"$`, tc.validCode)
	ctx.Step(`^the annual budget is -([0-9,]+)$`, tc.negativeBudget)

	// Given steps - Validation scenarios
	ctx.Step(`^I attempt to create departments with these codes:$`, tc.attemptCreateWithCodes)
	ctx.Step(`^I attempt to create departments with these budgets:$`, tc.attemptCreateWithBudgets)
	ctx.Step(`^I attempt to create departments with these statuses:$`, tc.attemptCreateWithStatuses)

	// Given steps - Failure scenarios
	ctx.Step(`^the database connection is lost$`, tc.databaseLost)
	ctx.Step(`^Redis is unavailable$`, tc.redisUnavailable)
	ctx.Step(`^a department with code "([^"]*)" exists in database$`, tc.existsInDatabase)
	ctx.Step(`^(\d+) concurrent requests to create department "([^"]*)"$`, tc.concurrentRequests)
	ctx.Step(`^departments "([^"]*)", "([^"]*)", and "([^"]*)" are accessed most frequently$`, tc.frequentlyAccessed)

	// When steps - HTTP requests
	ctx.Step(`^I send a POST request to "([^"]*)"$`, tc.sendPOST)
	ctx.Step(`^I send a POST request to "([^"]*)" with code "([^"]*)"$`, tc.sendPOSTWithCode)
	ctx.Step(`^I send a GET request to "([^"]*)"$`, tc.sendGET)
	ctx.Step(`^I send a GET request to "([^"]*)" with valid data$`, tc.sendGETWithData)
	ctx.Step(`^I send a PUT request to "([^"]*)" with:$`, tc.sendPUTWith)
	ctx.Step(`^I send a PUT request to "([^"]*)" with status "([^"]*)"$`, tc.sendPUTWithStatus)
	ctx.Step(`^I send a PUT request to "([^"]*)" with budget \$([0-9,]+)$`, tc.sendPUTWithBudget)
	ctx.Step(`^I send a PUT request to "([^"]*)" with valid data$`, tc.sendPUTValidData)
	ctx.Step(`^I send a PUT request to "([^"]*)" with updated data$`, tc.sendPUTUpdated)
	ctx.Step(`^I send a DELETE request to "([^"]*)"$`, tc.sendDELETE)
	ctx.Step(`^I send an OPTIONS request to "([^"]*)"$`, tc.sendOPTIONS)
	ctx.Step(`^I send any request to "([^"]*)"$`, tc.sendAnyRequest)
	ctx.Step(`^I send a POST request to "([^"]*)" with invalid JSON$`, tc.sendPOSTInvalidJSON)

	// When steps - Actions
	ctx.Step(`^I update the budget to \$([0-9,]+)$`, tc.updateBudget)
	ctx.Step(`^I update the budget again to \$([0-9,]+)$`, tc.updateBudgetAgain)
	ctx.Step(`^the cache warming process runs$`, tc.cacheWarming)

	// Then steps - Response validation
	ctx.Step(`^the response status should be (\d+) ([^"]*)$`, tc.responseStatusShouldBe)
	ctx.Step(`^the response should contain the department with code "([^"]*)"$`, tc.responseContainsDepartment)
	ctx.Step(`^the department should have an ID assigned$`, tc.hasIDAssigned)
	ctx.Step(`^the created_at timestamp should be set$`, tc.createdAtSet)
	ctx.Step(`^the updated_at timestamp should be set$`, tc.updatedAtSet)
	ctx.Step(`^the error message should indicate "([^"]*)"$`, tc.errorMessageIndicates)
	ctx.Step(`^the response should contain department with code "([^"]*)"$`, tc.responseContainsDepartmentCode)
	ctx.Step(`^the response should contain department with ID (\d+)$`, tc.responseContainsDepartmentID)
	ctx.Step(`^the response should include all department fields$`, tc.responseHasAllFields)
	ctx.Step(`^the response should contain (\d+) departments$`, tc.responseContainsDepartmentsCount)
	ctx.Step(`^all departments should have complete data$`, tc.allDepartmentsComplete)
	ctx.Step(`^all returned departments should have status "([^"]*)"$`, tc.allHaveStatus)

	// Then steps - Update validation
	ctx.Step(`^the department should be updated with new values$`, tc.departmentUpdated)
	ctx.Step(`^the updated_at timestamp should be newer than created_at$`, tc.updatedAtNewer)
	ctx.Step(`^the department status should be "([^"]*)"$`, tc.statusShouldBe)
	ctx.Step(`^the department status should be "([^"]*)" by default$`, tc.statusDefault)

	// Then steps - Delete validation
	ctx.Step(`^the department should no longer exist in the system$`, tc.noLongerExists)

	// Then steps - Cache validation
	ctx.Step(`^the department should be stored in Redis cache$`, tc.storedInCache)
	ctx.Step(`^the cache TTL should be (\d+) hour$`, tc.cacheTTL)
	ctx.Step(`^the response should include "([^"]*)" header$`, tc.responseIncludesHeader)
	ctx.Step(`^the database should not be queried$`, tc.databaseNotQueried)
	ctx.Step(`^the database should be queried$`, tc.databaseQueried)
	ctx.Step(`^the result should be cached for future requests$`, tc.cachedForFuture)
	ctx.Step(`^the Redis cache for "([^"]*)" should be invalidated$`, tc.cacheInvalidated)
	ctx.Step(`^the next GET request should show "([^"]*)"$`, tc.nextGETShows)

	// Then steps - Budget validation
	ctx.Step(`^the budget history should contain one entry$`, tc.budgetHistoryOneEntry)
	ctx.Step(`^the entry should show initial budget of \$([0-9,]+)$`, tc.entryShowsInitialBudget)
	ctx.Step(`^the budget history should contain (\d+) entries$`, tc.budgetHistoryEntries)
	ctx.Step(`^the entries should show the progression: \$([0-9,]+) → \$([0-9,]+) → \$([0-9,]+)$`, tc.budgetProgression)
	ctx.Step(`^the utilization percentage should be (\d+)%$`, tc.utilizationPercentage)
	ctx.Step(`^the remaining budget should be \$([0-9,]+)$`, tc.remainingBudget)
	ctx.Step(`^the report should include total allocated budget$`, tc.reportIncludesTotalAllocated)
	ctx.Step(`^the report should include total utilization across departments$`, tc.reportIncludesTotalUtilization)
	ctx.Step(`^the report should list each department's budget status$`, tc.reportListsBudgetStatus)

	// Then steps - Validation results
	ctx.Step(`^only valid codes should be accepted$`, tc.onlyValidCodesAccepted)
	ctx.Step(`^invalid codes should return (\d+) Bad Request$`, tc.invalidCodesReturnBadRequest)
	ctx.Step(`^only valid budgets should be accepted$`, tc.onlyValidBudgetsAccepted)
	ctx.Step(`^only "([^"]*)" and "([^"]*)" should be accepted$`, tc.onlyStatusesAccepted)

	// Then steps - Concurrent requests
	ctx.Step(`^only (\d+) department should be created successfully$`, tc.onlyOneDepartmentCreated)
	ctx.Step(`^(\d+) requests should fail with (\d+) Conflict$`, tc.requestsFailWithConflict)

	// Then steps - Cache warming
	ctx.Step(`^these departments should be pre-loaded into Redis cache$`, tc.preloadedIntoCache)
	ctx.Step(`^subsequent requests should show "([^"]*)"$`, tc.subsequentRequestsShow)

	// Then steps - API standards
	ctx.Step(`^the response should have "([^"]*)" header$`, tc.responseHasHeader)
	ctx.Step(`^the response should include appropriate CORS headers$`, tc.responseHasCORSHeaders)
	ctx.Step(`^the following endpoints should be available:$`, tc.endpointsAvailable)

	// Scenario hooks
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		tc.reset()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, nil
	})
}

// Stub implementations - these will be implemented when wiring up with actual handlers

func (tc *DepartmentTestContext) apiIsRunning() error                                   { return nil }
func (tc *DepartmentTestContext) redisIsAvailable() error                               { return nil }
func (tc *DepartmentTestContext) databaseIsInitialized() error                          { return nil }
func (tc *DepartmentTestContext) haveDepartmentData(code, name string) error            { return nil }
func (tc *DepartmentTestContext) setAnnualBudget(budget string) error                   { return nil }
func (tc *DepartmentTestContext) setDepartmentHead(head string) error                   { return nil }
func (tc *DepartmentTestContext) setStatus(status string) error                         { return nil }
func (tc *DepartmentTestContext) noStatusSpecified() error                              { return nil }
func (tc *DepartmentTestContext) departmentExists(code string) error                    { return nil }
func (tc *DepartmentTestContext) departmentWithIDExists(id int) error                   { return nil }
func (tc *DepartmentTestContext) departmentAlreadyExists(code string) error             { return nil }
func (tc *DepartmentTestContext) departmentWithStatusExists(code, status string) error  { return nil }
func (tc *DepartmentTestContext) cacheEmpty(code string) error                          { return nil }
func (tc *DepartmentTestContext) cachedInRedis(code string) error                       { return nil }
func (tc *DepartmentTestContext) notInCache(code string) error                          { return nil }
func (tc *DepartmentTestContext) cachedMinutesAgo(code string, minutes int) error       { return nil }
func (tc *DepartmentTestContext) multipleDepartmentsExist(table *godog.Table) error     { return nil }
func (tc *DepartmentTestContext) differentStatuses() error                              { return nil }
func (tc *DepartmentTestContext) differentBudgets() error                               { return nil }
func (tc *DepartmentTestContext) createDepartmentWithBudget(code, budget string) error  { return nil }
func (tc *DepartmentTestContext) departmentHasBudget(code, budget string) error         { return nil }
func (tc *DepartmentTestContext) departmentSpent(amount string) error                   { return nil }
func (tc *DepartmentTestContext) lowercaseCode(code string) error                       { return nil }
func (tc *DepartmentTestContext) codeWithLetters(code string, count int) error          { return nil }
func (tc *DepartmentTestContext) validCode(code string) error                           { return nil }
func (tc *DepartmentTestContext) negativeBudget(budget string) error                    { return nil }
func (tc *DepartmentTestContext) attemptCreateWithCodes(table *godog.Table) error       { return nil }
func (tc *DepartmentTestContext) attemptCreateWithBudgets(table *godog.Table) error     { return nil }
func (tc *DepartmentTestContext) attemptCreateWithStatuses(table *godog.Table) error    { return nil }
func (tc *DepartmentTestContext) databaseLost() error                                   { return nil }
func (tc *DepartmentTestContext) redisUnavailable() error                               { return nil }
func (tc *DepartmentTestContext) existsInDatabase(code string) error                    { return nil }
func (tc *DepartmentTestContext) concurrentRequests(count int, code string) error       { return nil }
func (tc *DepartmentTestContext) frequentlyAccessed(code1, code2, code3 string) error   { return nil }
func (tc *DepartmentTestContext) sendPOST(endpoint string) error                        { return nil }
func (tc *DepartmentTestContext) sendPOSTWithCode(endpoint, code string) error          { return nil }
func (tc *DepartmentTestContext) sendGET(endpoint string) error                         { return nil }
func (tc *DepartmentTestContext) sendGETWithData(endpoint string) error                 { return nil }
func (tc *DepartmentTestContext) sendPUTWith(endpoint string, table *godog.Table) error { return nil }
func (tc *DepartmentTestContext) sendPUTWithStatus(endpoint, status string) error       { return nil }
func (tc *DepartmentTestContext) sendPUTWithBudget(endpoint, budget string) error       { return nil }
func (tc *DepartmentTestContext) sendPUTValidData(endpoint string) error                { return nil }
func (tc *DepartmentTestContext) sendPUTUpdated(endpoint string) error                  { return nil }
func (tc *DepartmentTestContext) sendDELETE(endpoint string) error                      { return nil }
func (tc *DepartmentTestContext) sendOPTIONS(endpoint string) error                     { return nil }
func (tc *DepartmentTestContext) sendAnyRequest(endpoint string) error                  { return nil }
func (tc *DepartmentTestContext) sendPOSTInvalidJSON(endpoint string) error             { return nil }
func (tc *DepartmentTestContext) updateBudget(budget string) error                      { return nil }
func (tc *DepartmentTestContext) updateBudgetAgain(budget string) error                 { return nil }
func (tc *DepartmentTestContext) cacheWarming() error                                   { return nil }
func (tc *DepartmentTestContext) responseStatusShouldBe(code int, desc string) error    { return nil }
func (tc *DepartmentTestContext) responseContainsDepartment(code string) error          { return nil }
func (tc *DepartmentTestContext) hasIDAssigned() error                                  { return nil }
func (tc *DepartmentTestContext) createdAtSet() error                                   { return nil }
func (tc *DepartmentTestContext) updatedAtSet() error                                   { return nil }
func (tc *DepartmentTestContext) errorMessageIndicates(message string) error            { return nil }
func (tc *DepartmentTestContext) responseContainsDepartmentCode(code string) error      { return nil }
func (tc *DepartmentTestContext) responseContainsDepartmentID(id int) error             { return nil }
func (tc *DepartmentTestContext) responseHasAllFields() error                           { return nil }
func (tc *DepartmentTestContext) responseContainsDepartmentsCount(count int) error      { return nil }
func (tc *DepartmentTestContext) allDepartmentsComplete() error                         { return nil }
func (tc *DepartmentTestContext) allHaveStatus(status string) error                     { return nil }
func (tc *DepartmentTestContext) departmentUpdated() error                              { return nil }
func (tc *DepartmentTestContext) updatedAtNewer() error                                 { return nil }
func (tc *DepartmentTestContext) statusShouldBe(status string) error                    { return nil }
func (tc *DepartmentTestContext) statusDefault(status string) error                     { return nil }
func (tc *DepartmentTestContext) noLongerExists() error                                 { return nil }
func (tc *DepartmentTestContext) storedInCache() error                                  { return nil }
func (tc *DepartmentTestContext) cacheTTL(hours int) error                              { return nil }
func (tc *DepartmentTestContext) responseIncludesHeader(header string) error            { return nil }
func (tc *DepartmentTestContext) databaseNotQueried() error                             { return nil }
func (tc *DepartmentTestContext) databaseQueried() error                                { return nil }
func (tc *DepartmentTestContext) cachedForFuture() error                                { return nil }
func (tc *DepartmentTestContext) cacheInvalidated(code string) error                    { return nil }
func (tc *DepartmentTestContext) nextGETShows(status string) error                      { return nil }
func (tc *DepartmentTestContext) budgetHistoryOneEntry() error                          { return nil }
func (tc *DepartmentTestContext) entryShowsInitialBudget(budget string) error           { return nil }
func (tc *DepartmentTestContext) budgetHistoryEntries(count int) error                  { return nil }
func (tc *DepartmentTestContext) budgetProgression(b1, b2, b3 string) error             { return nil }
func (tc *DepartmentTestContext) utilizationPercentage(pct int) error                   { return nil }
func (tc *DepartmentTestContext) remainingBudget(amount string) error                   { return nil }
func (tc *DepartmentTestContext) reportIncludesTotalAllocated() error                   { return nil }
func (tc *DepartmentTestContext) reportIncludesTotalUtilization() error                 { return nil }
func (tc *DepartmentTestContext) reportListsBudgetStatus() error                        { return nil }
func (tc *DepartmentTestContext) onlyValidCodesAccepted() error                         { return nil }
func (tc *DepartmentTestContext) invalidCodesReturnBadRequest(code int) error           { return nil }
func (tc *DepartmentTestContext) onlyValidBudgetsAccepted() error                       { return nil }
func (tc *DepartmentTestContext) onlyStatusesAccepted(s1, s2 string) error              { return nil }
func (tc *DepartmentTestContext) onlyOneDepartmentCreated(count int) error              { return nil }
func (tc *DepartmentTestContext) requestsFailWithConflict(count, code int) error        { return nil }
func (tc *DepartmentTestContext) preloadedIntoCache() error                             { return nil }
func (tc *DepartmentTestContext) subsequentRequestsShow(status string) error            { return nil }
func (tc *DepartmentTestContext) responseHasHeader(header string) error                 { return nil }
func (tc *DepartmentTestContext) responseHasCORSHeaders() error                         { return nil }
func (tc *DepartmentTestContext) endpointsAvailable(table *godog.Table) error           { return nil }

func (tc *DepartmentTestContext) reset() {
	tc.data = make(map[string]interface{})
}
