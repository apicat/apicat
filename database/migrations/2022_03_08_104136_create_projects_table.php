<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('projects', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('user_id')->comment('项目管理者id');
            $table->string('icon')->nullable()->comment('项目图标');
            $table->string('name')->comment('项目名称');
            $table->tinyInteger('visibility')->comment('项目可见:0私有,1公开');
            $table->string('description')->nullable()->comment('项目描述');
            $table->timestamps();

            $table->softDeletes();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('projects');
    }
};
