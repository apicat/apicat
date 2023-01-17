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
        Schema::create('iterations', function (Blueprint $table) {
            $table->collation = 'utf8mb4_general_ci';

            $table->id();
            $table->bigInteger('project_id')->index()->comment('项目id');
            $table->bigInteger('user_id')->comment('迭代创建者id');
            $table->string('title')->comment('迭代名称');
            $table->string('description')->nullable()->comment('迭代描述');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('iterations');
    }
};
