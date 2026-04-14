/**
 * Admin API barrel export
 * Centralized exports for all admin API modules
 */

import dashboardAPI from './dashboard'
import usersAPI from './users'
import redeemAPI from './redeem'
import promoAPI from './promo'
import announcementsAPI from './announcements'
import settingsAPI from './settings'
import systemAPI from './system'
import subscriptionsAPI from './subscriptions'
import userAttributesAPI from './userAttributes'
import dataManagementAPI from './dataManagement'
import apiKeysAPI from './apiKeys'
import backupAPI from './backup'
import payAPI from './pay'
import proxyAdminAPI from './proxy'

/**
 * Unified admin API object for convenient access
 */
export const adminAPI = {
  dashboard: dashboardAPI,
  users: usersAPI,
  redeem: redeemAPI,
  promo: promoAPI,
  announcements: announcementsAPI,
  settings: settingsAPI,
  system: systemAPI,
  subscriptions: subscriptionsAPI,
  userAttributes: userAttributesAPI,
  dataManagement: dataManagementAPI,
  apiKeys: apiKeysAPI,
  backup: backupAPI,
  pay: payAPI,
  proxy: proxyAdminAPI,
}

export {
  dashboardAPI,
  usersAPI,
  redeemAPI,
  promoAPI,
  announcementsAPI,
  settingsAPI,
  systemAPI,
  subscriptionsAPI,
  userAttributesAPI,
  dataManagementAPI,
  apiKeysAPI,
  backupAPI,
  payAPI,
  proxyAdminAPI,
}

export default adminAPI

// Re-export types used by components
export type { BalanceHistoryItem } from './users'
export type { BackupAgentHealth, DataManagementConfig } from './dataManagement'
