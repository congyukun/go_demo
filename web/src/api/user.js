import request from '@/utils/request'

// 获取用户列表 - GET /api/v1/users
export function getUsers(params) {
  return request({
    url: '/users',
    method: 'get',
    params: {
      page: params.page || 1,
      size: params.page_size || params.size || 10
    }
  })
}

// 获取用户详情 - GET /api/v1/users/:id
export function getUserById(id) {
  return request({
    url: `/users/${id}`,
    method: 'get'
  })
}

// 创建用户 - POST /api/v1/users
export function createUser(data) {
  return request({
    url: '/users',
    method: 'post',
    data
  })
}

// 更新用户 - PUT /api/v1/users/:id
export function updateUser(id, data) {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data
  })
}

// 删除用户 - DELETE /api/v1/users/:id
export function deleteUser(id) {
  return request({
    url: `/users/${id}`,
    method: 'delete'
  })
}

// 获取用户统计 - GET /api/v1/users/stats
export function getUserStats() {
  return request({
    url: '/users/stats',
    method: 'get'
  })
}

// 更新当前用户资料 - PUT /api/v1/users/profile
export function updateProfile(data) {
  return request({
    url: '/users/profile',
    method: 'put',
    data
  })
}

// 修改密码 - PUT /api/v1/users/password
export function changepassword(data) {
  return request({
    url: '/users/password',
    method: 'put',
    data
  })
}
