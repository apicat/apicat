<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\Hash;
use App\Models\User;

class AdminPassword extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'admin:password';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '重置系统管理员密码';

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
        if (!$admin = User::find(1)) {
            return $this->error('管理员账号不存在');
        }

        $data = [];
        $data['password'] = $this->ask('请输入新的管理员密码');

        $validator = Validator::make($data, [
            'password' => 'required|string|min:8|max:255',
        ]);

        if ($validator->fails()) {
            return $this->error($validator->errors()->first());
        }

        $admin->password = Hash::make($data['password']);
        $admin->save();

        $this->info('修改成功');
        return 0;
    }
}
