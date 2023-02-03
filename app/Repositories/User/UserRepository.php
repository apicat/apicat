<?php

namespace App\Repositories\User;

use Illuminate\Contracts\Auth\Authenticatable;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;
use App\Models\ProjectMember;
use App\Models\User;

class UserRepository
{
    /**
     * 修改我的账户信息
     * @param array $data 要修改的内容
     * @return boolean 成功true 失败false
     */
    public static function editMyAccountInfo(array $data)
    {
        // 不能修改密码
        if (isset($data['password'])) {
            unset($data['password']);
        }

        if (!$data) {
            return false;
        }

        foreach ($data as $k => $v) {
            Auth::user()->$k = $v;
        }

        return Auth::user()->save();
    }

    /**
     * 修改用户的账户信息
     * @param int $userID 用户id
     * @param array $data 修改信息
     * @return bool
     */
    public static function editUserAccountInfo(int $userID, array $data)
    {
        if (!$data) {
            return false;
        }

        $saveData = [];

        if (isset($data['email'])) {
            $saveData['email'] = $data['email'];
        }
        if (isset($data['name'])) {
            $saveData['name'] = $data['name'];
        }
        if (isset($data['avatar'])) {
            $saveData['avatar'] = $data['avatar'];
        }
        if (isset($data['authority'])) {
            $saveData['authority'] = $data['authority'];
        }
        if (isset($data['password'])) {
            $saveData['password'] = Hash::make($data['password']);
        }
        if (isset($data['invitation_token'])) {
            $saveData['invitation_token'] = $data['invitation_token'];
        }

        if (!$saveData) {
            return false;
        }

        return (bool)User::where('id', $userID)->update($saveData);
    }

    /**
     * 校验我的密码
     * @param string $password
     * @return bool 正确true 错误false
     */
    public static function checkMyPassword($password)
    {
        return Hash::check($password, Auth::user()->password);
    }

    /**
     * 修改我的密码
     * @param string $password
     * @return boolean 成功true 失败false
     */
    public static function changeMyPassword($password)
    {
        Auth::user()->password = Hash::make($password);
        return (bool)Auth::user()->save();
    }

    /**
     * 检查邮箱除指定用户外是否已存在
     * @param string $email
     * @param int $userID
     * @return bool true: 邮箱已被其他用户使用  false: 邮箱未被其他用户使用
     */
    public static function checkEmailExists(string $email, int $userID)
    {
        $originUserID = User::where('email', $email)->value('id');
        if (!$originUserID) {
            return false;
        }

        return $originUserID !== $userID;
    }

    /**
     * 获取用户人数
     * @return mixed 用户数量
     */
    public static function userCount()
    {
        return User::count();
    }

    /**
     * 获取所有用户
     * @param int $offset 游标起始位置
     * @param int $limit 获取数量
     * @return mixed user实例
     */
    public static function users(int $offset = 1, int $limit = 15)
    {
        return User::offset($offset)->limit($limit)->latest()->get();
    }

    /**
     * 查询用户信息
     * @param int $userID 用户id
     * @param array|string $select 查询字段
     * @return User 用户信息
     */
    public static function getUserInfo(int $userID, array|string $select = [])
    {
        if ($select) {
            if (is_string($select)) {
                return User::where('id', $userID)->value($select);
            } else {
                return User::select($select)->find($userID);
            }
        }
        return User::find($userID);
    }

    /**
     * 判断用户是否存在
     * @param int $userID
     * @return bool true: 存在  false: 不存在
     */
    public static function checkUserExists(int $userID)
    {
        return (bool)User::where('id', $userID)->value('id');
    }

    /**
     * 检查用户是否有团队管理权限
     * @param Authenticatable|User|int $user 用户对象或用户id
     * @return boolean true是 false不是
     */
    public static function hasAuthority($user)
    {
        if (is_int($user)) {
            $authority = User::where('id', $user)->value('authority');
            if ($authority === null) {
                return false;
            }
        } else {
            $authority = $user->authority;
        }

        return 2 > $authority;
    }

    /**
     * 检查用户是否为团队管理员
     * @param Authenticatable|User|int $user 用户对象或用户id
     * @return boolean true是 false不是
     */
    public static function isAdmin($user)
    {
        if (is_int($user)) {
            return 1 === User::where('id', $user)->value('authority');
        }

        return 1 === $user->authority;
    }

    /**
     * 检查用户是否为团队超级管理员
     * @param Authenticatable|User|int $user 用户对象或用户id
     * @return boolean true是 false不是
     */
    public static function isSuperAdmin($user)
    {
        if (is_int($user)) {
            return 0 === User::where('id', $user)->value('authority');
        }

        return 0 === $user->authority;
    }

    /**
     * 删除团队成员
     * @param int $userID 成员id
     * @return boolean true成功 false失败
     */
    public static function remove(int $userID)
    {
        $user = User::find($userID);

        if (!$user) {
            return false;
        }

        try {
            DB::transaction(function () use ($user) {
                $user->delete();

                ProjectMember::where('user_id', $user->id)->delete();
            });
        } catch (\Exception $e) {
            return false;
        }

        return true;
    }

    /**
     * 添加用户
     * @param array $data
     * @return User|false 成功返回用户实例, 失败返回false
     */
    public static function addAccount(array $data)
    {
        $user = User::create([
            'email' => $data['email'],
            'name' => explode('@', $data['email'])[0],
            'avatar' => $data['avatar'] ?? null,
            'password' => Hash::make($data['password']),
            'authority' => $data['authority'],
            'invitation_token' => $data['invitation_token'] ?? null,
        ]);

        return $user ? $user : false;
    }

    /**
     * 查询多用户信息
     * @param array $userIDs 用户id数组
     */
    public static function getUsersByUserID(array $userIDs)
    {
        return User::whereIn('id', $userIDs)->get();
    }

    /**
     * 获取所有用户id
     * @return array
     */
    public static function userIds()
    {
        return User::pluck('id')->toArray();
    }

    /**
     * 通过id获取用户名称
     *
     * @param int $id 用户id
     * @param boolean $deleted 是否包含已删除用户
     * @return string
     */
    public static function name($id, $deleted = false)
    {
        if ($deleted) {
            return User::withTrashed()->where('id', $id)->value('name');
        }

        return User::where('id', $id)->value('name');
    }

    /**
     * 获取所有用户id和用户名称键值对
     *
     * @param boolean $deleted 是否包含已删除用户
     * @return \Illuminate\Support\Collection
     */
    public static function idNameArr($deleted = false)
    {
        if ($deleted) {
            return User::withTrashed()->pluck('name', 'id');
        }

        return User::pluck('name', 'id');
    }
}
