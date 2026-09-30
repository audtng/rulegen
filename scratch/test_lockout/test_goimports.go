package rules

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

type LockoutManager interface {
	CheckLockout(userID string) error
	IsLocked(userID string) bool
	RecordFailedAttempt(userID string) error
	ResetAttempts(userID string) error
}

type LockoutPolicy struct {
	MaxPasswordAttempts int
	MaxOTPAttempts      int
	MaxFailedAttempts   int
}

func testVulnDiscardRateLimiter(limiter *rate.Limiter) {
	// ruleid: auth-lockout-bypass
	_ = limiter.Allow()
}

func testVulnRateLimiterMissingReturn(w http.ResponseWriter, limiter *rate.Limiter) {
	// ruleid: auth-lockout-bypass
	if !limiter.Allow() {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
	}
}

func testVulnDiscardLockoutCheck(lm LockoutManager, userID string) {
	// ruleid: auth-lockout-bypass
	_ = lm.CheckLockout(userID)
}

func testVulnLockoutMissingReturn(w http.ResponseWriter, lm LockoutManager, userID string) {
	// ruleid: auth-lockout-bypass
	if lm.IsLocked(userID) {
		http.Error(w, "account is locked out", http.StatusUnauthorized)
	}
}

func testVulnDiscardRecordFailedAttempt(lm LockoutManager, userID string) {
	// ruleid: auth-lockout-bypass
	_ = lm.RecordFailedAttempt(userID)
}

func testVulnZeroOTPAttemptsPolicy() *LockoutPolicy {
	return &LockoutPolicy{
		MaxPasswordAttempts: 5,
		// ruleid: auth-lockout-bypass
		MaxOTPAttempts: 0,
	}
}

func testSafeRateLimiter(w http.ResponseWriter, limiter *rate.Limiter) {
	// ok: auth-lockout-bypass
	if !limiter.Allow() {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}
}

func testSafeLockoutCheck(w http.ResponseWriter, lm LockoutManager, userID string) {
	// ok: auth-lockout-bypass
	if err := lm.CheckLockout(userID); err != nil {
		http.Error(w, "account locked out", http.StatusUnauthorized)
		return
	}
}

func testSafeIsLocked(w http.ResponseWriter, lm LockoutManager, userID string) {
	// ok: auth-lockout-bypass
	if lm.IsLocked(userID) {
		http.Error(w, "account locked out", http.StatusUnauthorized)
		return
	}
}

func testSafeRecordFailedAttempt(lm LockoutManager, userID string) error {
	// ok: auth-lockout-bypass
	if err := lm.RecordFailedAttempt(userID); err != nil {
		return err
	}
	return nil
}

func testSafeRateLimiterWait(ctx context.Context, w http.ResponseWriter, limiter *rate.Limiter) {
	// ok: auth-lockout-bypass
	if err := limiter.Wait(ctx); err != nil {
		http.Error(w, "rate limit wait failed", http.StatusTooManyRequests)
		return
	}
}

func testSafePolicy() *LockoutPolicy {
	return &LockoutPolicy{
		MaxPasswordAttempts: 5,
		// ok: auth-lockout-bypass
		MaxOTPAttempts: 5,
	}
}
