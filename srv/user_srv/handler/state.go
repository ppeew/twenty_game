package handler

//
//var mutex sync.Mutex
//
//// 查询用户状态（先在redis访问，找不到再去mysql，mysql找到后插入redis）
//func (s *UserServer) GetUserState(ctx context.Context, req *user.UserIDInfo) (*user.UserStateResponse, error) {
//	result, err := global.RedisDB.Get(ctx, NameState(req.Id)).Result()
//	if err == redis.Nil {
//		//redis没找到,访问mysql
//		var u model.User
//		//加互斥锁（防止缓存击穿）
//		mutex.Lock()
//		//再次查询
//		res, err := global.RedisDB.Get(ctx, NameState(req.Id)).Result()
//		if err == nil {
//			//找到redis
//			state, _ := strconv.Atoi(res)
//			return &user.UserStateResponse{State: uint32(state)}, nil
//		}
//		if res := global.MysqlDB.Where("id=?", req.Id).First(&u); res.RowsAffected == 0 {
//			//mysql没找到(为防止缓存穿透，设置空值，维持30s)
//			global.RedisDB.Set(ctx, NameState(req.Id), "", 30*time.Second)
//			return &user.UserStateResponse{}, errors.New("无该用户")
//		}
//		//设置到随机过期时间缓存(防止缓存雪崩问题)
//		global.RedisDB.Set(ctx, NameState(req.Id), u.UserState, time.Duration(rand.Intn(100))*time.Minute)
//		mutex.Unlock()
//		return &user.UserStateResponse{State: u.UserState}, nil
//	}
//	//找到了直接返回（要判断是不是空值）
//	if result == "" {
//		return &user.UserStateResponse{}, errors.New("无该用户")
//	}
//	state, _ := strconv.Atoi(result)
//	return &user.UserStateResponse{State: uint32(state)}, nil
//}
//
//// 修改用户状态（修改mysql，删除redis）
//func (s *UserServer) UpdateUserState(ctx context.Context, req *user.UpdateUserStateInfo) (*emptypb.Empty, error) {
//	res := global.MysqlDB.Model(&model.User{}).Where("id=?", req.Id).Update("user_state", req.State)
//	if res.Error != nil {
//		//更新失败
//		return &emptypb.Empty{}, errors.New("更新状态失败")
//	}
//	//删除redis
//	global.RedisDB.Del(ctx, NameState(req.Id))
//	return &emptypb.Empty{}, nil
//}
//
//func NameState(id uint32) string {
//	return fmt.Sprintf("User:userState:%d", id)
//}
