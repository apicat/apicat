<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;

class AdminCreate extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'admin:create';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '创建系统管理员';

    /**
     * 用户表名称
     *
     * @var string
     */
    protected $table = 'users';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        if (DB::table($this->table)->find(1)) {
            return $this->error('管理员账号已经存在，无法重复创建管理员。');
        }

        $data = [];

        $data['email'] = $this->ask('请输入管理员邮箱');

        $validator = Validator::make($data, [
            'email' => 'required|email|max:255|unique:' . $this->table,
        ]);

        if ($validator->fails()) {
            return $this->error($validator->errors()->first());
        }

        $data['password'] = $this->ask('请输入管理员密码');

        $validator = Validator::make($data, [
            'password' => 'required|string|min:8|max:255',
        ]);

        if ($validator->fails()) {
            return $this->error($validator->errors()->first());
        }

        $data['password'] = Hash::make($data['password']);
        $data['name'] = explode('@', $data['email'])[0];
        $data['authority'] = 0;
        $data['id'] = 1;

        if (!DB::table($this->table)->insert($data)) {
            return $this->error('系统管理员创建失败');
        }

        $this->info('创建成功');
        return 0;
    }
}
